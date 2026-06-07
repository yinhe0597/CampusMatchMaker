package service

import (
	"context"
	"errors"
	"fmt"

	"campus_collab/internal/domain/classrepo"
	"campus_collab/internal/infra/cache"
	"campus_collab/internal/service/dto"
	appErr "campus_collab/pkg/errors"
	"campus_collab/pkg/utils"

	"gorm.io/gorm"
)

// TimetableInheritor 课表继承接口（避免循环依赖）
type TimetableInheritor interface {
	InheritTimetable(ctx context.Context, userID, classID uint) (int, error)
}

// ClassService 班级服务
type ClassService struct {
	classRepo        classrepo.ClassRepository
	ttInheritor      TimetableInheritor
	cache            *cache.Cache
}

// NewClassService 创建班级服务
func NewClassService(cr classrepo.ClassRepository, tti TimetableInheritor, cacheStore *cache.Cache) *ClassService {
	return &ClassService{
		classRepo:   cr,
		ttInheritor: tti,
		cache:       cacheStore,
	}
}

// CreateClass 创建班级
func (s *ClassService) CreateClass(ctx context.Context, userID uint, req *dto.CreateClassRequest) (*dto.CreateClassResult, error) {
	inviteCode := utils.GenerateInviteCode()

	class := &classrepo.Class{
		SchoolID:      req.SchoolID,
		Grade:         req.Grade,
		Department:    req.Department,
		Name:          req.Name,
		CreatorUserID: userID,
		InviteCode:    &inviteCode,
	}

	if err := s.classRepo.CreateClass(ctx, class); err != nil {
		return nil, appErr.Wrap(1000, "创建班级失败", err)
	}

	// 创建者自动成为 owner 成员
	member := &classrepo.ClassMember{
		ClassID: class.ID,
		UserID:  userID,
		Role:    classrepo.RoleOwner,
		Status:  classrepo.MemberActive,
	}
	if err := s.classRepo.CreateMember(ctx, member); err != nil {
		return nil, appErr.Wrap(1000, "添加成员失败", err)
	}

	return &dto.CreateClassResult{
		ID:              class.ID,
		Name:            class.Name,
		InviteCode:      inviteCode,
		TimetableStatus: class.TimetableStatus,
	}, nil
}

// GetClassDetail 获取班级详情
func (s *ClassService) GetClassDetail(ctx context.Context, userID, classID uint) (*dto.ClassDetailResult, error) {
	// 验证用户是成员
	isMember, err := s.classRepo.IsMember(ctx, classID, userID)
	if err != nil {
		return nil, appErr.New(1000, "查询失败")
	}
	if !isMember {
		return nil, appErr.ErrNoPermission
	}

	// 尝试从缓存读取
	cacheKey := fmt.Sprintf("class:detail:%d", classID)
	var cached dto.ClassDetailResult
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	class, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrClassNotFound
		}
		return nil, appErr.Wrap(1000, "查询班级失败", err)
	}

	// 获取角色决定是否返回邀请码
	role, _ := s.classRepo.GetMemberRole(ctx, classID, userID)
	var inviteCode *string
	if role == classrepo.RoleOwner || role == classrepo.RoleAdmin {
		inviteCode = class.InviteCode
	}

	memberCount, _ := s.classRepo.CountMembers(ctx, classID)

	result := &dto.ClassDetailResult{
		ID:              class.ID,
		SchoolID:        class.SchoolID,
		Grade:           class.Grade,
		Department:      class.Department,
		Name:            class.Name,
		Code:            class.Code,
		CreatorUserID:   class.CreatorUserID,
		InviteCode:      inviteCode,
		TimetableStatus: class.TimetableStatus,
		MemberCount:     memberCount,
		MyRole:          role,
		CreatedAt:       class.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	// 写入缓存
	_ = s.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// JoinClass 加入班级（根据 invite_code）
func (s *ClassService) JoinClass(ctx context.Context, userID, classID uint, req *dto.JoinClassRequest) (*dto.JoinClassResult, error) {
	class, err := s.classRepo.GetClassByID(ctx, classID)
	if err != nil {
		return nil, appErr.ErrClassNotFound
	}

	// 验证邀请码
	if class.InviteCode == nil || *class.InviteCode != req.InviteCode {
		return nil, appErr.ErrInviteCode
	}

	// 检查是否已是成员
	existing, err := s.classRepo.GetMember(ctx, classID, userID)
	if err == nil && existing.Status == classrepo.MemberActive {
		return nil, appErr.ErrAlreadyMember
	}

	if existing != nil && existing.Status == classrepo.MemberLeft {
		// 重新激活
		if err := s.classRepo.UpdateMemberStatus(ctx, classID, userID, classrepo.MemberActive); err != nil {
			return nil, appErr.Wrap(1000, "加入班级失败", err)
		}
	} else {
		// 新成员
		member := &classrepo.ClassMember{
			ClassID: classID,
			UserID:  userID,
			Role:    classrepo.RoleMember,
			Status:  classrepo.MemberActive,
		}
		if err := s.classRepo.CreateMember(ctx, member); err != nil {
			return nil, appErr.Wrap(1000, "加入班级失败", err)
		}
	}

	// 如果班级已有课表，触发继承
	if class.TimetableStatus == 1 && s.ttInheritor != nil {
		s.ttInheritor.InheritTimetable(ctx, userID, classID)
	}

	return &dto.JoinClassResult{
		ClassID:    class.ID,
		ClassName:  class.Name,
		InviteCode: *class.InviteCode,
		Role:       classrepo.RoleMember,
	}, nil
}

// ListMembers 获取成员列表
func (s *ClassService) ListMembers(ctx context.Context, userID, classID uint, page, pageSize int) (*dto.MemberListResult, error) {
	isMember, err := s.classRepo.IsMember(ctx, classID, userID)
	if err != nil || !isMember {
		return nil, appErr.ErrNoPermission
	}

	members, total, err := s.classRepo.ListMembers(ctx, classID, page, pageSize)
	if err != nil {
		return nil, appErr.Wrap(1000, "查询成员失败", err)
	}

	items := make([]dto.MemberItem, 0, len(members))
	for _, m := range members {
		items = append(items, dto.MemberItem{
			UserID:   m.UserID,
			Nickname: m.Nickname,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	return &dto.MemberListResult{
		Members:  items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// RemoveMember 移除成员
func (s *ClassService) RemoveMember(ctx context.Context, operatorID, classID, targetUserID uint) error {
	// 验证操作者身份
	role, err := s.classRepo.GetMemberRole(ctx, classID, operatorID)
	if err != nil {
		return appErr.ErrNoPermission
	}
	if role != classrepo.RoleOwner && role != classrepo.RoleAdmin {
		return appErr.ErrNoPermission
	}

	// 不能移除 owner
	targetRole, err := s.classRepo.GetMemberRole(ctx, classID, targetUserID)
	if err != nil {
		return appErr.New(1002, "目标成员不存在")
	}
	if targetRole == classrepo.RoleOwner {
		return appErr.ErrNoPermission
	}

	return s.classRepo.UpdateMemberStatus(ctx, classID, targetUserID, classrepo.MemberLeft)
}

// LookupClassByCode 通过邀请码查找班级基本信息
func (s *ClassService) LookupClassByCode(ctx context.Context, inviteCode string) (*dto.LookupClassResult, error) {
	class, err := s.classRepo.GetClassByInviteCode(ctx, inviteCode)
	if err != nil {
		return nil, appErr.ErrClassNotFound
	}

	count, _ := s.classRepo.CountMembers(ctx, class.ID)

	return &dto.LookupClassResult{
		ID:          class.ID,
		Name:        class.Name,
		Grade:       class.Grade,
		MemberCount: count,
	}, nil
}

// ListMyClasses 获取用户加入的班级列表
func (s *ClassService) ListMyClasses(ctx context.Context, userID uint) ([]*dto.MyClassItem, error) {
	classes, _, err := s.classRepo.ListUserClasses(ctx, userID)
	if err != nil {
		return nil, appErr.Wrap(1000, "查询班级列表失败", err)
	}

	items := make([]*dto.MyClassItem, 0, len(classes))
	for _, c := range classes {
		role, _ := s.classRepo.GetMemberRole(ctx, c.ID, userID)
		count, _ := s.classRepo.CountMembers(ctx, c.ID)

		item := &dto.MyClassItem{
			ID:              c.ID,
			Name:            c.Name,
			Grade:           c.Grade,
			Role:            role,
			MemberCount:     count,
			TimetableStatus: c.TimetableStatus,
			CreatedAt:       c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if c.Department != nil {
			item.Department = *c.Department
		}
		if c.InviteCode != nil {
			item.InviteCode = *c.InviteCode
		}
		items = append(items, item)
	}
	return items, nil
}

// IsMemberOfClass 检查用户是否为班级成员（供其他 service 调用）
func (s *ClassService) IsMemberOfClass(ctx context.Context, classID, userID uint) (bool, error) {
	return s.classRepo.IsMember(ctx, classID, userID)
}

// GetMemberIDs 获取班级所有活跃成员ID（供投票模块使用）
func (s *ClassService) GetMemberIDs(ctx context.Context, classID uint) ([]uint, error) {
	return s.classRepo.GetScopeMemberIDs(ctx, classID)
}
