package service

import (
	"context"
	"errors"
	"time"

	"campus_collab/internal/domain/classrepo"
	"campus_collab/internal/domain/timetable"
	"campus_collab/internal/service/dto"
	appErr "campus_collab/pkg/errors"

	"gorm.io/gorm"
)

// TimetableService 课表服务
type TimetableService struct {
	ttRepo    timetable.TimetableRepository
	classRepo classrepo.ClassRepository
}

// NewTimetableService 创建课表服务
func NewTimetableService(ttRepo timetable.TimetableRepository, cr classrepo.ClassRepository) *TimetableService {
	return &TimetableService{ttRepo: ttRepo, classRepo: cr}
}

// ===== 班级公共课表 =====

// CreateClassTimetable 录入班级公共课表
func (s *TimetableService) CreateClassTimetable(ctx context.Context, userID, classID uint, req *dto.CreateClassTimetableRequest) (*dto.CreateClassTimetableResult, error) {
	// 验证是班级成员
	if ok, _ := s.classRepo.IsMember(ctx, classID, userID); !ok {
		return nil, appErr.ErrNoPermission
	}

	// 检查是否已有课表
	existing, _ := s.ttRepo.GetClassTimetables(ctx, classID)
	if len(existing) > 0 {
		return nil, appErr.ErrTimetableExists
	}

	// 构建条目
	entries := make([]*timetable.ClassTimetable, 0, len(req.Entries))
	for _, e := range req.Entries {
		if !timetable.CheckPeriodRange(e.PeriodStart, e.PeriodEnd) {
			return nil, appErr.ErrPeriodInvalid
		}
		entries = append(entries, &timetable.ClassTimetable{
			ClassID:           classID,
			DayOfWeek:         e.DayOfWeek,
			PeriodStart:       e.PeriodStart,
			PeriodEnd:         e.PeriodEnd,
			CourseName:        e.CourseName,
			Teacher:           e.Teacher,
			Room:              e.Room,
			ContributorUserID: &userID,
			Version:           1,
		})
	}

	created, err := s.ttRepo.CreateClassTimetables(ctx, entries)
	if err != nil {
		return nil, appErr.Wrap(1000, "录入课表失败", err)
	}

	return &dto.CreateClassTimetableResult{
		CreatedCount:    created,
		TimetableStatus: 1,
	}, nil
}

// GetClassTimetable 获取班级公共课表
func (s *TimetableService) GetClassTimetable(ctx context.Context, userID, classID uint) (*dto.ClassTimetableResult, error) {
	if ok, _ := s.classRepo.IsMember(ctx, classID, userID); !ok {
		return nil, appErr.ErrNoPermission
	}

	entries, err := s.ttRepo.GetClassTimetables(ctx, classID)
	if err != nil {
		return nil, appErr.Wrap(1000, "查询课表失败", err)
	}

	items := make([]dto.ClassTimetableItem, 0, len(entries))
	for _, e := range entries {
		items = append(items, dto.ClassTimetableItem{
			ID:          e.ID,
			DayOfWeek:   e.DayOfWeek,
			PeriodStart: e.PeriodStart,
			PeriodEnd:   e.PeriodEnd,
			CourseName:  e.CourseName,
			Teacher:     e.Teacher,
			Room:        e.Room,
			Version:     e.Version,
		})
	}

	return &dto.ClassTimetableResult{
		ClassID:      classID,
		Entries:      items,
		TotalEntries: len(items),
	}, nil
}

// UpdateClassTimetable 更新班级公共课表（全量替换）
func (s *TimetableService) UpdateClassTimetable(ctx context.Context, userID, classID uint, req *dto.UpdateClassTimetableRequest) (*dto.UpdateClassTimetableResult, error) {
	// 验证权限（owner/admin）
	role, err := s.classRepo.GetMemberRole(ctx, classID, userID)
	if err != nil || (role != classrepo.RoleOwner && role != classrepo.RoleAdmin) {
		return nil, appErr.ErrNoPermission
	}

	// 验证节次范围
	for _, e := range req.Entries {
		if !timetable.CheckPeriodRange(e.PeriodStart, e.PeriodEnd) {
			return nil, appErr.ErrPeriodInvalid
		}
	}

	// 软删除旧条目
	if err := s.ttRepo.DeleteClassTimetables(ctx, classID); err != nil {
		return nil, appErr.Wrap(1000, "更新课表失败", err)
	}

	// 创建新条目
	entries := make([]*timetable.ClassTimetable, 0, len(req.Entries))
	for _, e := range req.Entries {
		entries = append(entries, &timetable.ClassTimetable{
			ClassID:           classID,
			DayOfWeek:         e.DayOfWeek,
			PeriodStart:       e.PeriodStart,
			PeriodEnd:         e.PeriodEnd,
			CourseName:        e.CourseName,
			Teacher:           e.Teacher,
			Room:              e.Room,
			ContributorUserID: &userID,
			Version:           1,
		})
	}

	created, err := s.ttRepo.CreateClassTimetables(ctx, entries)
	if err != nil {
		return nil, appErr.Wrap(1000, "更新课表失败", err)
	}

	// 清除所有成员的继承课表，重新继承
	s.ttRepo.DeletePersonalByClass(ctx, classID, timetable.PTInherited)
	memberIDs, _ := s.classRepo.GetScopeMemberIDs(ctx, classID)
	affectedCount := int64(0)
	for _, mid := range memberIDs {
		cnt, _ := s.inheritForUser(ctx, mid, classID)
		affectedCount += int64(cnt)
	}

	return &dto.UpdateClassTimetableResult{
		UpdatedCount:    created,
		AffectedMembers: affectedCount,
	}, nil
}

// ===== 个人课表 =====

// CreatePersonalTimetable 添加个人课表条目
func (s *TimetableService) CreatePersonalTimetable(ctx context.Context, userID uint, req *dto.CreatePersonalTimetableRequest) (*dto.CreatePersonalTimetableResult, error) {
	if !timetable.CheckPeriodRange(req.PeriodStart, req.PeriodEnd) {
		return nil, appErr.ErrPeriodInvalid
	}

	// 检查时间冲突
	conflict, err := s.ttRepo.CheckTimeConflict(ctx, userID, req.ClassID, req.DayOfWeek, req.PeriodStart, req.PeriodEnd, 0)
	if err != nil {
		return nil, appErr.Wrap(1000, "添加课程失败", err)
	}
	if conflict {
		return nil, appErr.ErrTimeConflict
	}

	entry := &timetable.PersonalTimetable{
		UserID:      userID,
		ClassID:     &req.ClassID,
		DayOfWeek:   req.DayOfWeek,
		PeriodStart: req.PeriodStart,
		PeriodEnd:   req.PeriodEnd,
		CourseName:  req.CourseName,
		Source:      timetable.PTPersonal,
	}

	if err := s.ttRepo.CreatePersonalTimetable(ctx, entry); err != nil {
		return nil, appErr.Wrap(1000, "添加课程失败", err)
	}

	return &dto.CreatePersonalTimetableResult{
		ID:     entry.ID,
		Source: timetable.PTPersonal,
	}, nil
}

// GetPersonalTimetable 获取个人完整课表（继承+个人）
func (s *TimetableService) GetPersonalTimetable(ctx context.Context, userID, classID uint) (*dto.PersonalTimetableResult, error) {
	if ok, _ := s.classRepo.IsMember(ctx, classID, userID); !ok {
		return nil, appErr.ErrNoPermission
	}

	entries, err := s.ttRepo.GetPersonalTimetable(ctx, userID, classID)
	if err != nil {
		return nil, appErr.Wrap(1000, "查询课表失败", err)
	}

	// 如果用户还没有继承课表但班级有公共课表，自动继承
	hasInherited, _ := s.ttRepo.HasInherited(ctx, userID, classID)
	if !hasInherited {
		s.inheritForUser(ctx, userID, classID)
		entries, _ = s.ttRepo.GetPersonalTimetable(ctx, userID, classID)
	}

	items := make([]dto.PersonalTimetableItem, 0, len(entries))
	inheritedCount := 0
	personalCount := 0
	for _, e := range entries {
		if e.Source == timetable.PTInherited {
			inheritedCount++
		} else {
			personalCount++
		}
		items = append(items, dto.PersonalTimetableItem{
			ID:           e.ID,
			DayOfWeek:    e.DayOfWeek,
			PeriodStart:  e.PeriodStart,
			PeriodEnd:    e.PeriodEnd,
			CourseName:   e.CourseName,
			Source:       e.Source,
			IsOverridden: e.IsOverridden == 1,
		})
	}

	return &dto.PersonalTimetableResult{
		UserID:         userID,
		ClassID:        classID,
		Entries:        items,
		InheritedCount: inheritedCount,
		PersonalCount:  personalCount,
	}, nil
}

// UpdatePersonalTimetable 修改个人课表条目
func (s *TimetableService) UpdatePersonalTimetable(ctx context.Context, userID, entryID uint, req *dto.UpdatePersonalTimetableRequest) error {
	entry, err := s.ttRepo.GetPersonalTimetableByID(ctx, entryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErr.ErrTimetableNotFound
		}
		return appErr.Wrap(1000, "查询失败", err)
	}
	if entry.UserID != userID {
		return appErr.ErrNoPermission
	}

	if req.CourseName != nil {
		entry.CourseName = *req.CourseName
	}
	if req.IsOverridden != nil {
		if *req.IsOverridden {
			entry.IsOverridden = 1
		} else {
			entry.IsOverridden = 0
		}
	}

	return s.ttRepo.UpdatePersonalTimetable(ctx, entry)
}

// DeletePersonalTimetable 删除个人课表条目
func (s *TimetableService) DeletePersonalTimetable(ctx context.Context, userID, entryID uint) error {
	entry, err := s.ttRepo.GetPersonalTimetableByID(ctx, entryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErr.ErrTimetableNotFound
		}
		return appErr.Wrap(1000, "查询失败", err)
	}
	if entry.UserID != userID {
		return appErr.ErrNoPermission
	}

	return s.ttRepo.SoftDeletePersonalTimetable(ctx, entryID)
}

// ===== 继承 =====

// InheritTimetable 实现 TimetableInheritor 接口（供 ClassService 调用）
func (s *TimetableService) InheritTimetable(ctx context.Context, userID, classID uint) (int, error) {
	return s.inheritForUser(ctx, userID, classID)
}

func (s *TimetableService) inheritForUser(ctx context.Context, userID, classID uint) (int, error) {
	classEntries, err := s.ttRepo.GetClassTimetables(ctx, classID)
	if err != nil || len(classEntries) == 0 {
		return 0, nil
	}

	personal := make([]*timetable.PersonalTimetable, 0, len(classEntries))
	for _, ce := range classEntries {
		personal = append(personal, &timetable.PersonalTimetable{
			UserID:              userID,
			ClassID:             &classID,
			DayOfWeek:           ce.DayOfWeek,
			PeriodStart:         ce.PeriodStart,
			PeriodEnd:           ce.PeriodEnd,
			CourseName:          ce.CourseName,
			Source:              timetable.PTInherited,
			RefClassTimetableID: &ce.ID,
		})
	}

	created, err := s.ttRepo.CreatePersonalTimetables(ctx, personal)
	if err != nil {
		return 0, err
	}
	return created, nil
}

// ===== 纠错 =====

func (s *TimetableService) CreateCorrection(ctx context.Context, userID uint, req *dto.CreateCorrectionRequest) (*dto.CreateCorrectionResult, error) {
	// 验证公共课表条目存在
	_, err := s.ttRepo.GetClassTimetableByID(ctx, req.ClassTimetableID)
	if err != nil {
		return nil, appErr.ErrTimetableNotFound
	}

	correction := &timetable.TimetableCorrection{
		ClassTimetableID:      req.ClassTimetableID,
		ReporterUserID:        userID,
		CorrectionType:        req.CorrectionType,
		Description:           req.Description,
		SuggestedCourseName:   req.SuggestedCourseName,
		SuggestedPeriodStart:  req.SuggestedPeriodStart,
		SuggestedPeriodEnd:    req.SuggestedPeriodEnd,
		Status:                timetable.TCStatusPending,
	}

	if err := s.ttRepo.CreateCorrection(ctx, correction); err != nil {
		return nil, appErr.Wrap(1000, "提交纠错失败", err)
	}

	return &dto.CreateCorrectionResult{
		ID:      correction.ID,
		Status:  timetable.TCStatusPending,
		Message: "纠错已提交，等待审核",
	}, nil
}

func (s *TimetableService) ListCorrections(ctx context.Context, userID, classID uint, status int, page, pageSize int) (*dto.CorrectionListResult, error) {
	// 验证权限（owner/admin）
	role, err := s.classRepo.GetMemberRole(ctx, classID, userID)
	if err != nil || (role != classrepo.RoleOwner && role != classrepo.RoleAdmin) {
		return nil, appErr.ErrNoPermission
	}

	corrections, total, err := s.ttRepo.ListCorrections(ctx, classID, status, page, pageSize)
	if err != nil {
		return nil, appErr.Wrap(1000, "查询纠错列表失败", err)
	}

	items := make([]dto.CorrectionItem, 0, len(corrections))
	for _, c := range corrections {
		item := dto.CorrectionItem{
			ID:               c.ID,
			ClassTimetableID: c.ClassTimetableID,
			ReporterUserID:   c.ReporterUserID,
			CorrectionType:   c.CorrectionType,
			Status:           c.Status,
			CreatedAt:        c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if c.Description != nil {
			item.Description = *c.Description
		}
		if c.SuggestedCourseName != nil {
			item.SuggestedCourseName = *c.SuggestedCourseName
		}
		item.SuggestedPeriodStart = c.SuggestedPeriodStart
		item.SuggestedPeriodEnd = c.SuggestedPeriodEnd
		items = append(items, item)
	}

	return &dto.CorrectionListResult{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *TimetableService) ReviewCorrection(ctx context.Context, userID, correctionID uint, req *dto.ReviewCorrectionRequest) error {
	correction, err := s.ttRepo.GetCorrectionByID(ctx, correctionID)
	if err != nil {
		return appErr.ErrCorrectionNotFound
	}
	if correction.Status != timetable.TCStatusPending {
		return appErr.ErrCorrectionHandled
	}

	// 获取公共课条目的班级ID
	ctEntry, err := s.ttRepo.GetClassTimetableByID(ctx, correction.ClassTimetableID)
	if err != nil {
		return appErr.ErrTimetableNotFound
	}
	classID := ctEntry.ClassID

	// 验证权限（owner/admin）
	role, err := s.classRepo.GetMemberRole(ctx, classID, userID)
	if err != nil || (role != classrepo.RoleOwner && role != classrepo.RoleAdmin) {
		return appErr.ErrNoPermission
	}

	now := time.Now()
	correction.ReviewedBy = &userID
	correction.ReviewedAt = &now
	correction.ResolvedAt = &now

	if req.Action == "approve" {
		correction.Status = timetable.TCStatusApproved

		// 更新公共课表
		if correction.SuggestedCourseName != nil {
			ctEntry.CourseName = *correction.SuggestedCourseName
		}
		if correction.SuggestedPeriodStart != nil {
			ctEntry.PeriodStart = *correction.SuggestedPeriodStart
		}
		if correction.SuggestedPeriodEnd != nil {
			ctEntry.PeriodEnd = *correction.SuggestedPeriodEnd
		}
		ctEntry.Version++
		if err := s.ttRepo.UpdateClassTimetable(ctx, ctEntry); err != nil {
			return appErr.Wrap(1000, "更新课表失败", err)
		}
	} else {
		correction.Status = timetable.TCStatusRejected
	}

	return s.ttRepo.UpdateCorrection(ctx, correction)
}
