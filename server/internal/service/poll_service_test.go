package service

import (
	"context"
	"testing"
	"time"

	"campus_collab/internal/domain/classrepo"
	"campus_collab/internal/domain/poll"
	"campus_collab/internal/domain/timetable"
	"campus_collab/internal/engine/timeslot"
	"campus_collab/internal/service/dto"
)

// ===== Mock Repositories =====

type mockPollRepo struct {
	polls    map[uint]*poll.Poll
	options  map[uint][]*poll.PollOption
	votes    []*poll.PollVote
	results  map[uint]*poll.PollResult
	nextID   uint
}

func newMockPollRepo() *mockPollRepo {
	return &mockPollRepo{
		polls:   make(map[uint]*poll.Poll),
		options: make(map[uint][]*poll.PollOption),
		results: make(map[uint]*poll.PollResult),
		nextID:  1,
	}
}

func (m *mockPollRepo) CreatePoll(ctx context.Context, p *poll.Poll) error {
	p.ID = m.nextID
	m.nextID++
	m.polls[p.ID] = p
	return nil
}
func (m *mockPollRepo) GetPollByID(ctx context.Context, id uint) (*poll.Poll, error) {
	p, ok := m.polls[id]
	if !ok {
		return nil, context.DeadlineExceeded // simulate not found
	}
	return p, nil
}
func (m *mockPollRepo) ListPolls(ctx context.Context, scopeType string, scopeID uint, page, pageSize int) ([]*poll.Poll, int64, error) {
	var result []*poll.Poll
	for _, p := range m.polls {
		if p.ScopeType == scopeType && p.ScopeID == scopeID {
			result = append(result, p)
		}
	}
	return result, int64(len(result)), nil
}
func (m *mockPollRepo) UpdatePoll(ctx context.Context, p *poll.Poll) error {
	m.polls[p.ID] = p
	return nil
}
func (m *mockPollRepo) CreateOptions(ctx context.Context, opts []*poll.PollOption) (int, error) {
	for _, o := range opts {
		o.ID = m.nextID
		m.nextID++
		m.options[o.PollID] = append(m.options[o.PollID], o)
	}
	return len(opts), nil
}
func (m *mockPollRepo) GetOptionsByPoll(ctx context.Context, pollID uint) ([]*poll.PollOption, error) {
	return m.options[pollID], nil
}
func (m *mockPollRepo) GetOptionByID(ctx context.Context, id uint) (*poll.PollOption, error) {
	for _, opts := range m.options {
		for _, o := range opts {
			if o.ID == id {
				return o, nil
			}
		}
	}
	return nil, context.DeadlineExceeded
}
func (m *mockPollRepo) DeleteOptions(ctx context.Context, pollID uint) error {
	delete(m.options, pollID)
	return nil
}
func (m *mockPollRepo) CreateVote(ctx context.Context, v *poll.PollVote) error {
	v.ID = m.nextID
	m.nextID++
	m.votes = append(m.votes, v)
	return nil
}
func (m *mockPollRepo) CreateVotes(ctx context.Context, votes []*poll.PollVote) (int, error) {
	for _, v := range votes {
		v.ID = m.nextID
		m.nextID++
		m.votes = append(m.votes, v)
	}
	return len(votes), nil
}
func (m *mockPollRepo) UpdateVote(ctx context.Context, v *poll.PollVote) error {
	for i, existing := range m.votes {
		if existing.ID == v.ID {
			m.votes[i] = v
			return nil
		}
	}
	return nil
}
func (m *mockPollRepo) GetVoteByUnique(ctx context.Context, pollID, optionID, voterUserID uint) (*poll.PollVote, error) {
	for _, v := range m.votes {
		if v.PollID == pollID && v.OptionID == optionID && v.VoterUserID == voterUserID {
			return v, nil
		}
	}
	return nil, context.DeadlineExceeded // simulate not found
}
func (m *mockPollRepo) ListVotesByPoll(ctx context.Context, pollID uint) ([]*poll.PollVote, error) {
	var result []*poll.PollVote
	for _, v := range m.votes {
		if v.PollID == pollID {
			result = append(result, v)
		}
	}
	return result, nil
}
func (m *mockPollRepo) ListVotesByUserAndPoll(ctx context.Context, pollID, voterUserID uint) ([]*poll.PollVote, error) {
	var result []*poll.PollVote
	for _, v := range m.votes {
		if v.PollID == pollID && v.VoterUserID == voterUserID {
			result = append(result, v)
		}
	}
	return result, nil
}
func (m *mockPollRepo) CountVotersByPoll(ctx context.Context, pollID uint) (int64, error) {
	seen := make(map[uint]struct{})
	for _, v := range m.votes {
		if v.PollID == pollID {
			seen[v.VoterUserID] = struct{}{}
		}
	}
	return int64(len(seen)), nil
}
func (m *mockPollRepo) UpsertResult(ctx context.Context, result *poll.PollResult) error {
	key := result.PollID*1000 + result.OptionID
	m.results[key] = result
	return nil
}
func (m *mockPollRepo) GetResultsByPoll(ctx context.Context, pollID uint) ([]*poll.PollResult, error) {
	var result []*poll.PollResult
	for _, r := range m.results {
		if r.PollID == pollID {
			result = append(result, r)
		}
	}
	return result, nil
}
func (m *mockPollRepo) DeleteResults(ctx context.Context, pollID uint) error {
	for k := range m.results {
		if k/1000 == pollID {
			delete(m.results, k)
		}
	}
	return nil
}

type mockClassRepo struct {
	members map[uint][]uint // classID -> []userID
	classes map[uint]*classrepo.Class
}

func newMockClassRepo() *mockClassRepo {
	return &mockClassRepo{
		members: make(map[uint][]uint),
		classes: make(map[uint]*classrepo.Class),
	}
}
func (m *mockClassRepo) CreateClass(ctx context.Context, c *classrepo.Class) error { return nil }
func (m *mockClassRepo) GetClassByID(ctx context.Context, id uint) (*classrepo.Class, error) {
	c, ok := m.classes[id]
	if !ok {
		return &classrepo.Class{ID: id, TimetableStatus: 0}, nil
	}
	return c, nil
}
func (m *mockClassRepo) GetClassByInviteCode(ctx context.Context, code string) (*classrepo.Class, error) {
	return nil, nil
}
func (m *mockClassRepo) UpdateClass(ctx context.Context, c *classrepo.Class) error {
	m.classes[c.ID] = c
	return nil
}
func (m *mockClassRepo) CreateMember(ctx context.Context, mem *classrepo.ClassMember) error {
	m.members[mem.ClassID] = append(m.members[mem.ClassID], mem.UserID)
	return nil
}
func (m *mockClassRepo) GetMember(ctx context.Context, classID, userID uint) (*classrepo.ClassMember, error) {
	return nil, nil
}
func (m *mockClassRepo) GetMemberRole(ctx context.Context, classID, userID uint) (string, error) {
	return classrepo.RoleOwner, nil
}
func (m *mockClassRepo) IsMember(ctx context.Context, classID, userID uint) (bool, error) {
	for _, uid := range m.members[classID] {
		if uid == userID {
			return true, nil
		}
	}
	return false, nil
}
func (m *mockClassRepo) ListMembers(ctx context.Context, classID uint, page, pageSize int) ([]classrepo.MemberInfo, int64, error) {
	return nil, 0, nil
}
func (m *mockClassRepo) UpdateMemberStatus(ctx context.Context, classID, userID uint, status int) error {
	return nil
}
func (m *mockClassRepo) CountMembers(ctx context.Context, classID uint) (int64, error) {
	return int64(len(m.members[classID])), nil
}
func (m *mockClassRepo) GetScopeMemberIDs(ctx context.Context, classID uint) ([]uint, error) {
	return m.members[classID], nil
}
func (m *mockClassRepo) ListUserClasses(ctx context.Context, userID uint) ([]*classrepo.Class, int64, error) {
	return nil, 0, nil
}

type mockTTRepo struct {
	personalTimetables map[uint][]*timetable.PersonalTimetable // userID -> entries
}

func newMockTTRepo() *mockTTRepo {
	return &mockTTRepo{
		personalTimetables: make(map[uint][]*timetable.PersonalTimetable),
	}
}
func (m *mockTTRepo) GetPersonalTimetable(ctx context.Context, userID, classID uint) ([]*timetable.PersonalTimetable, error) {
	var result []*timetable.PersonalTimetable
	for _, e := range m.personalTimetables[userID] {
		if e.ClassID != nil && *e.ClassID == classID && e.DeletedAt == nil {
			result = append(result, e)
		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// Stub methods for the rest of the TimetableRepository interface
func (m *mockTTRepo) CreateClassTimetables(ctx context.Context, entries []*timetable.ClassTimetable) (int, error) {
	return 0, nil
}
func (m *mockTTRepo) GetClassTimetables(ctx context.Context, classID uint) ([]*timetable.ClassTimetable, error) {
	return nil, nil
}
func (m *mockTTRepo) GetClassTimetableByID(ctx context.Context, id uint) (*timetable.ClassTimetable, error) {
	return nil, nil
}
func (m *mockTTRepo) UpdateClassTimetable(ctx context.Context, entry *timetable.ClassTimetable) error {
	return nil
}
func (m *mockTTRepo) DeleteClassTimetables(ctx context.Context, classID uint) error {
	return nil
}
func (m *mockTTRepo) SoftDeleteClassTimetable(ctx context.Context, id uint) error {
	return nil
}
func (m *mockTTRepo) IncrementVersion(ctx context.Context, classID uint) error {
	return nil
}
func (m *mockTTRepo) CreatePersonalTimetable(ctx context.Context, entry *timetable.PersonalTimetable) error {
	return nil
}
func (m *mockTTRepo) CreatePersonalTimetables(ctx context.Context, entries []*timetable.PersonalTimetable) (int, error) {
	return 0, nil
}
func (m *mockTTRepo) GetPersonalTimetableByID(ctx context.Context, id uint) (*timetable.PersonalTimetable, error) {
	return nil, nil
}
func (m *mockTTRepo) UpdatePersonalTimetable(ctx context.Context, entry *timetable.PersonalTimetable) error {
	return nil
}
func (m *mockTTRepo) SoftDeletePersonalTimetable(ctx context.Context, id uint) error {
	return nil
}
func (m *mockTTRepo) DeletePersonalByClass(ctx context.Context, classID uint, source string) error {
	return nil
}
func (m *mockTTRepo) CountPersonalByClass(ctx context.Context, classID uint) (int64, error) {
	return 0, nil
}
func (m *mockTTRepo) HasInherited(ctx context.Context, userID, classID uint) (bool, error) {
	return false, nil
}
func (m *mockTTRepo) CheckTimeConflict(ctx context.Context, userID, classID uint, dayOfWeek, periodStart, periodEnd int, excludeID uint) (bool, error) {
	return false, nil
}
func (m *mockTTRepo) CreateCorrection(ctx context.Context, correction *timetable.TimetableCorrection) error {
	return nil
}
func (m *mockTTRepo) ListCorrections(ctx context.Context, classID uint, status, page, pageSize int) ([]*timetable.TimetableCorrection, int64, error) {
	return nil, 0, nil
}
func (m *mockTTRepo) GetCorrectionByID(ctx context.Context, id uint) (*timetable.TimetableCorrection, error) {
	return nil, nil
}
func (m *mockTTRepo) UpdateCorrection(ctx context.Context, correction *timetable.TimetableCorrection) error {
	return nil
}

// ===== Tests =====

func newTestPollService() (*PollService, *mockPollRepo, *mockClassRepo, *mockTTRepo) {
	pollRepo := newMockPollRepo()
	classRepo := newMockClassRepo()
	ttRepo := newMockTTRepo()
	svc := NewPollService(pollRepo, classRepo, ttRepo)
	return svc, pollRepo, classRepo, ttRepo
}

// TestCreatePoll_WithoutRecommend 不启用推荐
func TestCreatePoll_WithoutRecommend(t *testing.T) {
	svc, _, classRepo, _ := newTestPollService()
	classRepo.members[1] = []uint{10} // class1 has member userID=10

	deadline := time.Now().Add(24 * time.Hour)
	result, err := svc.CreatePoll(context.Background(), 10, &dto.CreatePollRequest{
		Title:         "测试投票",
		Description:   nil,
		ScopeType:     "class",
		ScopeID:       1,
		Deadline:      &deadline,
		AutoRecommend: false,
	})
	if err != nil {
		t.Fatalf("CreatePoll 失败: %v", err)
	}
	if result.PollID == 0 {
		t.Error("PollID 不应为 0")
	}
	if result.Status != poll.PollStatusDraft {
		t.Errorf("新创建投票状态应为 draft，实际: %s", result.Status)
	}
	if result.OptionsCreated != 0 {
		t.Errorf("未启用推荐时应为 0，实际: %d", result.OptionsCreated)
	}
}

// TestCreatePoll_WithRecommend 启用自动推荐
func TestCreatePoll_WithRecommend(t *testing.T) {
	svc, _, classRepo, ttRepo := newTestPollService()
	classRepo.members[1] = []uint{10, 20} // 2 members

	// 用户10的课表：周一第1-2节(480-575)
	ttRepo.personalTimetables[10] = []*timetable.PersonalTimetable{
		{DayOfWeek: 1, PeriodStart: 1, PeriodEnd: 2, Source: timetable.PTInherited,
			CourseName: "数学"},
	}
	cid := uint(1)
	ttRepo.personalTimetables[10][0].ClassID = &cid

	// 用户20的课表：周一第3-4节(595-690)
	ttRepo.personalTimetables[20] = []*timetable.PersonalTimetable{
		{DayOfWeek: 1, PeriodStart: 3, PeriodEnd: 4, Source: timetable.PTInherited,
			CourseName: "英语"},
	}
	ttRepo.personalTimetables[20][0].ClassID = &cid

	deadline := time.Now().Add(24 * time.Hour)
	result, err := svc.CreatePoll(context.Background(), 10, &dto.CreatePollRequest{
		Title:         "测试投票+推荐",
		ScopeType:     "class",
		ScopeID:       1,
		Deadline:      &deadline,
		AutoRecommend: true,
		TimePreference: &dto.TimePref{
			DayStartHour:      8,
			DayEndHour:        22,
			MinDurationMin:    60,
			MaxRecommendations: 5,
		},
	})
	if err != nil {
		t.Fatalf("CreatePoll 失败: %v", err)
	}
	if result.OptionsCreated == 0 {
		t.Error("2个用户都有课表 → 应有推荐选项（发现引擎排序 bug 已修复）")
	}
	t.Logf("创建投票 ID=%d, 推荐选项数=%d", result.PollID, result.OptionsCreated)
}

// TestOpenPoll 开启投票
func TestOpenPoll(t *testing.T) {
	svc, pollRepo, _, _ := newTestPollService()
	pollRepo.polls[1] = &poll.Poll{ID: 1, CreatorUserID: 10, Status: poll.PollStatusDraft}

	err := svc.OpenPoll(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("OpenPoll 失败: %v", err)
	}
	if pollRepo.polls[1].Status != poll.PollStatusOpen {
		t.Errorf("状态应为 open，实际: %s", pollRepo.polls[1].Status)
	}
}

// TestOpenPoll_NotDraft 非草稿状态不可开启
func TestOpenPoll_NotDraft(t *testing.T) {
	svc, pollRepo, _, _ := newTestPollService()
	pollRepo.polls[1] = &poll.Poll{ID: 1, CreatorUserID: 10, Status: poll.PollStatusOpen}

	err := svc.OpenPoll(context.Background(), 10, 1)
	if err == nil {
		t.Error("已开启的投票不应再被开启")
	}
}

// TestClosePoll 关闭投票
func TestClosePoll(t *testing.T) {
	svc, pollRepo, _, _ := newTestPollService()
	pollRepo.polls[1] = &poll.Poll{ID: 1, CreatorUserID: 10, Status: poll.PollStatusOpen}

	err := svc.ClosePoll(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("ClosePoll 失败: %v", err)
	}
	if pollRepo.polls[1].Status != poll.PollStatusClosed {
		t.Errorf("状态应为 closed，实际: %s", pollRepo.polls[1].Status)
	}
}

// TestSubmitVote 提交投票
func TestSubmitVote(t *testing.T) {
	svc, pollRepo, classRepo, _ := newTestPollService()
	classRepo.members[1] = []uint{10, 20} // user 20 is also a member
	pollRepo.polls[1] = &poll.Poll{
		ID: 1, CreatorUserID: 10, Status: poll.PollStatusOpen,
		ScopeType: "class", ScopeID: 1,
	}
	pollRepo.options[1] = []*poll.PollOption{
		{ID: 100, PollID: 1},
		{ID: 101, PollID: 1},
	}

	result, err := svc.SubmitVote(context.Background(), 20, 1, &dto.SubmitVoteRequest{
		Votes: []dto.VoteInput{
			{OptionID: 100, Choice: poll.ChoiceYes},
			{OptionID: 101, Choice: poll.ChoiceMaybe},
		},
	})
	if err != nil {
		t.Fatalf("SubmitVote 失败: %v", err)
	}
	if result.VotedCount != 2 {
		t.Errorf("应为2票，实际: %d", result.VotedCount)
	}
}

// TestSubmitVote_NotOpen 非开启状态不可投票
func TestSubmitVote_NotOpen(t *testing.T) {
	svc, pollRepo, _, _ := newTestPollService()
	pollRepo.polls[1] = &poll.Poll{
		ID: 1, CreatorUserID: 10, Status: poll.PollStatusDraft,
		ScopeType: "class", ScopeID: 1,
	}

	_, err := svc.SubmitVote(context.Background(), 20, 1, &dto.SubmitVoteRequest{
		Votes: []dto.VoteInput{{OptionID: 100, Choice: poll.ChoiceYes}},
	})
	if err == nil {
		t.Error("草稿状态不应允许投票")
	}
}

// TestFinalizePoll 确认最终时段
func TestFinalizePoll(t *testing.T) {
	svc, pollRepo, _, _ := newTestPollService()
	pollRepo.polls[1] = &poll.Poll{ID: 1, CreatorUserID: 10, Status: poll.PollStatusClosed}
	pollRepo.options[1] = []*poll.PollOption{{ID: 100, PollID: 1}}

	err := svc.FinalizePoll(context.Background(), 10, 1, &dto.FinalizePollRequest{
		FinalOptionID: 100,
	})
	if err != nil {
		t.Fatalf("FinalizePoll 失败: %v", err)
	}
	if pollRepo.polls[1].Status != poll.PollStatusFinalized {
		t.Errorf("状态应为 finalized，实际: %s", pollRepo.polls[1].Status)
	}
}

// TestCalculateResults 计算结果
func TestCalculateResults(t *testing.T) {
	svc, pollRepo, classRepo, _ := newTestPollService()
	classRepo.members[1] = []uint{10, 20} // user 10 and 20 are members
	pollRepo.options[1] = []*poll.PollOption{
		{ID: 1, PollID: 1},
		{ID: 2, PollID: 1},
	}
	pollRepo.votes = []*poll.PollVote{
		{ID: 1, PollID: 1, OptionID: 1, VoterUserID: 10, Choice: poll.ChoiceYes},
		{ID: 2, PollID: 1, OptionID: 1, VoterUserID: 20, Choice: poll.ChoiceYes},
		{ID: 3, PollID: 1, OptionID: 2, VoterUserID: 10, Choice: poll.ChoiceNo},
	}

	pollRepo.polls[1] = &poll.Poll{
		ID: 1, CreatorUserID: 10, Status: poll.PollStatusClosed,
		ScopeType: "class", ScopeID: 1,
	}

	// 调用 GetResults 会触发 calculateResults
	result, err := svc.GetResults(context.Background(), 20, 1)
	if err != nil {
		t.Fatalf("GetResults 失败: %v", err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("应有2个选项结果，实际: %d", len(result.Items))
	}
	// option 1: 2人 yes
	opt1 := result.Items[0]
	if opt1.YesCount != 2 || opt1.NoCount != 0 || opt1.TotalVotes != 2 {
		t.Errorf("选项1结果错误: yes=%d no=%d total=%d", opt1.YesCount, opt1.NoCount, opt1.TotalVotes)
	}
	// option 2: 1人 no
	opt2 := result.Items[1]
	if opt2.YesCount != 0 || opt2.NoCount != 1 || opt2.TotalVotes != 1 {
		t.Errorf("选项2结果错误: yes=%d no=%d total=%d", opt2.YesCount, opt2.NoCount, opt2.TotalVotes)
	}
}

// TestEditPoll_NotCreator 非创建者不可编辑
func TestEditPoll_NotCreator(t *testing.T) {
	svc, pollRepo, _, _ := newTestPollService()
	pollRepo.polls[1] = &poll.Poll{ID: 1, CreatorUserID: 10, Status: poll.PollStatusDraft}

	err := svc.EditPoll(context.Background(), 20, 1, &dto.EditPollRequest{})
	if err == nil {
		t.Error("非创建者不应能编辑")
	}
}

// TestCalculateFreeSlots_EmptyMembers 空班级
func TestCalculateFreeSlots_EmptyMembers(t *testing.T) {
	svc, _, classRepo, _ := newTestPollService()
	classRepo.members[1] = []uint{} // 无成员

	_, err := svc.calculateFreeSlots(context.Background(), 1, nil)
	if err == nil {
		t.Error("空班级应返回错误")
	}
}

// Benchmark for engine
func BenchmarkCalculateFreeSlots(b *testing.B) {
	cfg := timeslot.DefaultConfig()
	schedules := make([]timeslot.UserSchedule, 30)
	for i := range schedules {
		schedules[i] = timeslot.UserSchedule{
			UserID: string(rune(i)),
			Slots: []timeslot.OccupiedSlot{
				{DayOfWeek: 1, StartMinutes: 480, EndMinutes: 600},
				{DayOfWeek: 3, StartMinutes: 780, EndMinutes: 900},
			},
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timeslot.CalculateFreeSlots(schedules, cfg)
	}
}
