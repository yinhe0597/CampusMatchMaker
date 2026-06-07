package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"campus_collab/internal/domain/classrepo"
	"campus_collab/internal/domain/poll"
	"campus_collab/internal/domain/timetable"
	"campus_collab/internal/engine/timeslot"
	"campus_collab/internal/service/dto"
	appErr "campus_collab/pkg/errors"
)

// PollService 投票服务
type PollService struct {
	pollRepo  poll.PollRepository
	classRepo classrepo.ClassRepository
	ttRepo    timetable.TimetableRepository
}

func NewPollService(
	pollRepo poll.PollRepository,
	classRepo classrepo.ClassRepository,
	ttRepo timetable.TimetableRepository,
) *PollService {
	return &PollService{
		pollRepo:  pollRepo,
		classRepo: classRepo,
		ttRepo:    ttRepo,
	}
}

// CreatePoll 创建投票（可选自动推荐空闲时段）
func (s *PollService) CreatePoll(ctx context.Context, userID uint, req *dto.CreatePollRequest) (*dto.CreatePollResult, error) {
	// 验证用户是班级成员
	if req.ScopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, req.ScopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证成员资格失败", err)
		}
		if !ok {
			return nil, appErr.ErrNoPermission
		}
	}

	// 创建投票（草稿状态）
	p := &poll.Poll{
		CreatorUserID:   userID,
		Title:           req.Title,
		Description:     req.Description,
		ScopeType:       req.ScopeType,
		ScopeID:         req.ScopeID,
		Status:          poll.PollStatusDraft,
		Deadline:        req.Deadline,
		MinParticipants: 2,
	}
	if err := s.pollRepo.CreatePoll(ctx, p); err != nil {
		return nil, appErr.Wrap(1000, "创建投票失败", err)
	}

	optCreated := 0

	// 自动推荐空闲时段
	if req.AutoRecommend && req.ScopeType == "class" {
		results, calcErr := s.calculateFreeSlots(ctx, req.ScopeID, req.TimePreference)
		if calcErr == nil && len(results) > 0 {
			now := time.Now()
			options := make([]*poll.PollOption, 0, len(results))
			for i, r := range results {
				slotDate, slotStart, slotEnd := freeSlotToDateTime(r, now)
				opt := &poll.PollOption{
					PollID:             p.ID,
					SlotDate:           slotDate,
					SlotStartTime:      slotStart,
					SlotEndTime:        slotEnd,
					DayOfWeek:          r.DayOfWeek,
					IsRecommended:      1,
					RecommendationRate: &r.Rate,
					SortOrder:          i + 1,
				}
				options = append(options, opt)
			}
			n, createErr := s.pollRepo.CreateOptions(ctx, options)
			if createErr == nil {
				optCreated = n
			}
		}
	}

	return &dto.CreatePollResult{
		PollID:         p.ID,
		Status:         p.Status,
		OptionsCreated: optCreated,
	}, nil
}

// 计算空闲时段
func (s *PollService) calculateFreeSlots(ctx context.Context, classID uint, pref *dto.TimePref) ([]timeslot.FreeSlotResult, error) {
	// 获取班级所有成员
	memberIDs, err := s.classRepo.GetScopeMemberIDs(ctx, classID)
	if err != nil {
		return nil, err
	}
	if len(memberIDs) == 0 {
		return nil, fmt.Errorf("班级没有成员")
	}

	// 配置引擎
	cfg := timeslot.DefaultConfig()
	if pref != nil {
		if pref.DayStartHour > 0 {
			cfg.DayStartMinutes = pref.DayStartHour * 60
		}
		if pref.DayEndHour > 0 {
			cfg.DayEndMinutes = pref.DayEndHour * 60
		}
	}

	// 获取每个成员的课表并转换为 OccupiedSlot
	var schedules []timeslot.UserSchedule
	for _, uid := range memberIDs {
		entries, qErr := s.ttRepo.GetPersonalTimetable(ctx, uid, classID)
		if qErr != nil {
			continue
		}

		slots := make([]timeslot.OccupiedSlot, 0, len(entries))
		for _, e := range entries {
			if e.DeletedAt != nil {
				continue
			}
			slots = append(slots, timeslot.OccupiedSlot{
				DayOfWeek:    e.DayOfWeek,
				StartMinutes: timetable.PeriodStartToMinutes(e.PeriodStart),
				EndMinutes:   timetable.PeriodEndToMinutes(e.PeriodEnd),
			})
		}

		schedules = append(schedules, timeslot.UserSchedule{
			UserID: fmt.Sprintf("%d", uid),
			Slots:  slots,
		})
	}

	if len(schedules) == 0 {
		return nil, fmt.Errorf("没有可用的课表数据")
	}

	// 调用引擎
	results := timeslot.CalculateFreeSlots(schedules, cfg)
	results = timeslot.FillTotalCount(results, len(schedules))

	// 过滤最短时长
	minDuration := 60
	if pref != nil && pref.MinDurationMin > 0 {
		minDuration = pref.MinDurationMin
	}
	var filtered []timeslot.FreeSlotResult
	for _, r := range results {
		if r.EndMinutes-r.StartMinutes >= minDuration {
			filtered = append(filtered, r)
		}
	}

	// 限制推荐数量
	maxRec := 5
	if pref != nil && pref.MaxRecommendations > 0 {
		maxRec = pref.MaxRecommendations
	}
	if len(filtered) > maxRec {
		filtered = filtered[:maxRec]
	}

	return filtered, nil
}

// freeSlotToDateTime 将空闲时段结果转换为实际日期和时间字符串
func freeSlotToDateTime(slot timeslot.FreeSlotResult, now time.Time) (date, startTime, endTime string) {
	daysUntil := (slot.DayOfWeek - int(now.Weekday())) % 7
	if daysUntil <= 0 {
		daysUntil += 7
	}
	targetDate := now.AddDate(0, 0, daysUntil)
	date = targetDate.Format("2006-01-02")

	startHour := slot.StartMinutes / 60
	startMin := slot.StartMinutes % 60
	endHour := slot.EndMinutes / 60
	endMin := slot.EndMinutes % 60

	startTime = fmt.Sprintf("%02d:%02d:00", startHour, startMin)
	endTime = fmt.Sprintf("%02d:%02d:00", endHour, endMin)
	return
}

// ListPolls 获取投票列表
func (s *PollService) ListPolls(ctx context.Context, userID, scopeID uint, scopeType string, page, pageSize int) (*dto.PollListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	if scopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, scopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证权限失败", err)
		}
		if !ok {
			return nil, appErr.ErrPollNoAccess
		}
	}

	polls, total, err := s.pollRepo.ListPolls(ctx, scopeType, scopeID, page, pageSize)
	if err != nil {
		return nil, appErr.Wrap(1000, "查询投票列表失败", err)
	}

	items := make([]*dto.PollListItem, 0, len(polls))
	for _, p := range polls {
		voterCount, _ := s.pollRepo.CountVotersByPoll(ctx, p.ID)
		opts, _ := s.pollRepo.GetOptionsByPoll(ctx, p.ID)

		items = append(items, &dto.PollListItem{
			ID:              p.ID,
			Title:           p.Title,
			CreatorUserID:   p.CreatorUserID,
			Status:          p.Status,
			Deadline:        p.Deadline,
			MinParticipants: p.MinParticipants,
			VoterCount:      voterCount,
			TotalOptions:    len(opts),
			CreatedAt:       p.CreatedAt,
		})
	}

	return &dto.PollListResult{Polls: items, TotalCount: total}, nil
}

// GetPollDetail 获取投票详情
func (s *PollService) GetPollDetail(ctx context.Context, userID, pollID uint) (*dto.PollDetailResult, error) {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, appErr.ErrPollNotFound
	}

	if p.ScopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, p.ScopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证权限失败", err)
		}
		if !ok {
			return nil, appErr.ErrPollNoAccess
		}
	}

	opts, err := s.pollRepo.GetOptionsByPoll(ctx, pollID)
	if err != nil {
		return nil, appErr.Wrap(1000, "获取选项失败", err)
	}

	myVotes, _ := s.pollRepo.ListVotesByUserAndPoll(ctx, pollID, userID)
	myVoteRecords := make([]*dto.VoteRecord, 0, len(myVotes))
	for _, v := range myVotes {
		myVoteRecords = append(myVoteRecords, &dto.VoteRecord{
			OptionID: v.OptionID,
			Choice:   v.Choice,
		})
	}

	results, _ := s.pollRepo.GetResultsByPoll(ctx, pollID)
	resultMap := make(map[uint]*poll.PollResult)
	for _, r := range results {
		resultMap[r.OptionID] = r
	}

	voterCount, _ := s.pollRepo.CountVotersByPoll(ctx, pollID)

	optionResults := make([]*dto.PollOptionResult, 0, len(opts))
	for _, opt := range opts {
		item := &dto.PollOptionItem{
			ID:                 opt.ID,
			SlotDate:           opt.SlotDate,
			SlotStartTime:      opt.SlotStartTime,
			SlotEndTime:        opt.SlotEndTime,
			DayOfWeek:          opt.DayOfWeek,
			IsRecommended:      opt.IsRecommended,
			RecommendationRate: opt.RecommendationRate,
			SortOrder:          opt.SortOrder,
		}

		r := resultMap[opt.ID]
		yesCount, noCount, maybeCount, totalVotes := 0, 0, 0, 0
		if r != nil {
			yesCount = r.YesCount
			noCount = r.NoCount
			maybeCount = r.MaybeCount
			totalVotes = r.TotalVotes
		}

		optionResults = append(optionResults, &dto.PollOptionResult{
			Option:     *item,
			YesCount:   yesCount,
			NoCount:    noCount,
			MaybeCount: maybeCount,
			TotalVotes: totalVotes,
		})
	}

	return &dto.PollDetailResult{
		ID:              p.ID,
		Title:           p.Title,
		Description:     p.Description,
		ScopeType:       p.ScopeType,
		ScopeID:         p.ScopeID,
		Status:          p.Status,
		Deadline:        p.Deadline,
		MinParticipants: p.MinParticipants,
		FinalOptionID:   p.FinalOptionID,
		CreatorUserID:   p.CreatorUserID,
		Options:         optionResults,
		MyVotes:         myVoteRecords,
		VoterCount:      voterCount,
		CreatedAt:       p.CreatedAt,
	}, nil
}

// EditPoll 编辑投票
func (s *PollService) EditPoll(ctx context.Context, userID, pollID uint, req *dto.EditPollRequest) error {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return appErr.ErrPollNotFound
	}
	if p.CreatorUserID != userID {
		return appErr.ErrOnlyCreator
	}
	if p.Status != poll.PollStatusDraft {
		return appErr.ErrPollNotOpen
	}

	if req.Title != nil {
		p.Title = *req.Title
	}
	if req.Description != nil {
		p.Description = req.Description
	}
	if req.Deadline != nil {
		p.Deadline = req.Deadline
	}
	return s.pollRepo.UpdatePoll(ctx, p)
}

// OpenPoll 开启投票
func (s *PollService) OpenPoll(ctx context.Context, userID, pollID uint) error {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return appErr.ErrPollNotFound
	}
	if p.CreatorUserID != userID {
		return appErr.ErrOnlyCreator
	}
	if p.Status != poll.PollStatusDraft {
		return appErr.ErrPollClosed
	}
	p.Status = poll.PollStatusOpen
	return s.pollRepo.UpdatePoll(ctx, p)
}

// ClosePoll 关闭投票
func (s *PollService) ClosePoll(ctx context.Context, userID, pollID uint) error {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return appErr.ErrPollNotFound
	}
	if p.CreatorUserID != userID {
		return appErr.ErrOnlyCreator
	}
	if p.Status != poll.PollStatusOpen {
		return appErr.ErrPollClosed
	}
	p.Status = poll.PollStatusClosed
	now := time.Now()
	p.ClosedAt = &now
	if err := s.pollRepo.UpdatePoll(ctx, p); err != nil {
		return appErr.Wrap(1000, "关闭投票失败", err)
	}
	return nil
}

// SubmitVote 提交投票
func (s *PollService) SubmitVote(ctx context.Context, userID, pollID uint, req *dto.SubmitVoteRequest) (*dto.SubmitVoteResult, error) {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, appErr.ErrPollNotFound
	}

	if p.Status != poll.PollStatusOpen {
		return nil, appErr.ErrPollNotOpen
	}
	if p.Deadline != nil && time.Now().After(*p.Deadline) {
		return nil, appErr.ErrPollDeadline
	}

	if p.ScopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, p.ScopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证权限失败", err)
		}
		if !ok {
			return nil, appErr.ErrPollNoAccess
		}
	}

	opts, err := s.pollRepo.GetOptionsByPoll(ctx, pollID)
	if err != nil {
		return nil, appErr.Wrap(1000, "获取选项失败", err)
	}
	optMap := make(map[uint]*poll.PollOption)
	for _, o := range opts {
		optMap[o.ID] = o
	}

	for _, v := range req.Votes {
		if _, ok := optMap[v.OptionID]; !ok {
			return nil, appErr.ErrOptionNotBelong
		}
	}

	votedCount := 0
	for _, v := range req.Votes {
		existing, lookupErr := s.pollRepo.GetVoteByUnique(ctx, pollID, v.OptionID, userID)
		if lookupErr != nil {
			newVote := &poll.PollVote{
				PollID:      pollID,
				OptionID:    v.OptionID,
				VoterUserID: userID,
				Choice:      v.Choice,
			}
			if createErr := s.pollRepo.CreateVote(ctx, newVote); createErr != nil {
				continue
			}
		} else {
			existing.Choice = v.Choice
			if updateErr := s.pollRepo.UpdateVote(ctx, existing); updateErr != nil {
				continue
			}
		}
		votedCount++
	}

	go s.recalculateResults(context.Background(), pollID)

	return &dto.SubmitVoteResult{VotedCount: votedCount}, nil
}

// GetResults 获取投票结果
func (s *PollService) GetResults(ctx context.Context, userID, pollID uint) (*dto.PollResultsResult, error) {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, appErr.ErrPollNotFound
	}

	if p.ScopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, p.ScopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证权限失败", err)
		}
		if !ok {
			return nil, appErr.ErrPollNoAccess
		}
	}

	if err := s.calculateResults(ctx, pollID); err != nil {
		return nil, appErr.Wrap(1000, "计算结果失败", err)
	}

	opts, err := s.pollRepo.GetOptionsByPoll(ctx, pollID)
	if err != nil {
		return nil, appErr.Wrap(1000, "获取选项失败", err)
	}
	results, err := s.pollRepo.GetResultsByPoll(ctx, pollID)
	if err != nil {
		return nil, appErr.Wrap(1000, "获取结果失败", err)
	}

	resultMap := make(map[uint]*poll.PollResult)
	for _, r := range results {
		resultMap[r.OptionID] = r
	}

	items := make([]*dto.OptionResultItem, 0, len(opts))
	for _, opt := range opts {
		r := resultMap[opt.ID]
		yesCount, noCount, maybeCount, totalVotes := 0, 0, 0, 0
		if r != nil {
			yesCount = r.YesCount
			noCount = r.NoCount
			maybeCount = r.MaybeCount
			totalVotes = r.TotalVotes
		}
		items = append(items, &dto.OptionResultItem{
			OptionID:      opt.ID,
			SlotDate:      opt.SlotDate,
			SlotStartTime: opt.SlotStartTime,
			SlotEndTime:   opt.SlotEndTime,
			DayOfWeek:     opt.DayOfWeek,
			YesCount:      yesCount,
			NoCount:       noCount,
			MaybeCount:    maybeCount,
			TotalVotes:    totalVotes,
		})
	}

	return &dto.PollResultsResult{Items: items}, nil
}

// FinalizePoll 确认最终时段
func (s *PollService) FinalizePoll(ctx context.Context, userID, pollID uint, req *dto.FinalizePollRequest) error {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return appErr.ErrPollNotFound
	}
	if p.CreatorUserID != userID {
		return appErr.ErrOnlyCreator
	}
	if p.Status != poll.PollStatusClosed && p.Status != poll.PollStatusOpen {
		return appErr.ErrPollClosed
	}

	opt, err := s.pollRepo.GetOptionByID(ctx, req.FinalOptionID)
	if err != nil {
		return appErr.ErrOptionNotBelong
	}
	if opt.PollID != pollID {
		return appErr.ErrOptionNotBelong
	}

	p.FinalOptionID = &opt.ID
	p.Status = poll.PollStatusFinalized
	return s.pollRepo.UpdatePoll(ctx, p)
}

// calculateResults 计算投票结果
func (s *PollService) calculateResults(ctx context.Context, pollID uint) error {
	votes, err := s.pollRepo.ListVotesByPoll(ctx, pollID)
	if err != nil {
		return err
	}
	opts, err := s.pollRepo.GetOptionsByPoll(ctx, pollID)
	if err != nil {
		return err
	}

	type counts struct {
		yes, no, maybe int
	}
	optCounts := make(map[uint]*counts)
	for _, opt := range opts {
		optCounts[opt.ID] = &counts{}
	}
	for _, v := range votes {
		if c, ok := optCounts[v.OptionID]; ok {
			switch v.Choice {
			case poll.ChoiceYes:
				c.yes++
			case poll.ChoiceNo:
				c.no++
			case poll.ChoiceMaybe:
				c.maybe++
			}
		}
	}

	voterSet := make(map[uint]struct{})
	for _, v := range votes {
		voterSet[v.VoterUserID] = struct{}{}
	}
	voterCount := int64(len(voterSet))

	for _, opt := range opts {
		c := optCounts[opt.ID]
		totalVotes := c.yes + c.no + c.maybe
		rate := 0.0
		if voterCount > 0 {
			rate = math.Min(float64(totalVotes)/float64(voterCount), 1.0)
		}
		result := &poll.PollResult{
			PollID:            pollID,
			OptionID:          opt.ID,
			YesCount:          c.yes,
			NoCount:           c.no,
			MaybeCount:        c.maybe,
			TotalVotes:        totalVotes,
			ParticipationRate: rate,
			CalculatedAt:      time.Now(),
		}
		if err := s.pollRepo.UpsertResult(ctx, result); err != nil {
			return err
		}
	}

	return nil
}

// recalculateResults 异步重新计算
func (s *PollService) recalculateResults(ctx context.Context, pollID uint) {
	_ = s.calculateResults(ctx, pollID)
}

// GetMyVotes 获取用户的投票记录
func (s *PollService) GetMyVotes(ctx context.Context, userID, pollID uint) ([]*dto.VoteRecord, error) {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, appErr.ErrPollNotFound
	}
	if p.ScopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, p.ScopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证权限失败", err)
		}
		if !ok {
			return nil, appErr.ErrPollNoAccess
		}
	}

	votes, err := s.pollRepo.ListVotesByUserAndPoll(ctx, pollID, userID)
	if err != nil {
		return nil, appErr.Wrap(1000, "获取投票记录失败", err)
	}

	records := make([]*dto.VoteRecord, 0, len(votes))
	for _, v := range votes {
		records = append(records, &dto.VoteRecord{
			OptionID: v.OptionID,
			Choice:   v.Choice,
		})
	}
	return records, nil
}

// GetOptions 获取投票选项
func (s *PollService) GetOptions(ctx context.Context, userID, pollID uint) ([]*dto.PollOptionItem, error) {
	p, err := s.pollRepo.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, appErr.ErrPollNotFound
	}
	if p.ScopeType == "class" {
		ok, err := s.classRepo.IsMember(ctx, p.ScopeID, userID)
		if err != nil {
			return nil, appErr.Wrap(1000, "验证权限失败", err)
		}
		if !ok {
			return nil, appErr.ErrPollNoAccess
		}
	}

	opts, err := s.pollRepo.GetOptionsByPoll(ctx, pollID)
	if err != nil {
		return nil, appErr.Wrap(1000, "获取选项失败", err)
	}

	items := make([]*dto.PollOptionItem, 0, len(opts))
	for _, opt := range opts {
		items = append(items, &dto.PollOptionItem{
			ID:                 opt.ID,
			SlotDate:           opt.SlotDate,
			SlotStartTime:      opt.SlotStartTime,
			SlotEndTime:        opt.SlotEndTime,
			DayOfWeek:          opt.DayOfWeek,
			IsRecommended:      opt.IsRecommended,
			RecommendationRate: opt.RecommendationRate,
			SortOrder:          opt.SortOrder,
		})
	}
	return items, nil
}
