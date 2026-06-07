package poll

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// PollRepository 投票仓储接口
type PollRepository interface {
	// 投票
	CreatePoll(ctx context.Context, p *Poll) error
	GetPollByID(ctx context.Context, id uint) (*Poll, error)
	ListPolls(ctx context.Context, scopeType string, scopeID uint, page, pageSize int) ([]*Poll, int64, error)
	UpdatePoll(ctx context.Context, p *Poll) error

	// 选项
	CreateOptions(ctx context.Context, opts []*PollOption) (int, error)
	GetOptionsByPoll(ctx context.Context, pollID uint) ([]*PollOption, error)
	GetOptionByID(ctx context.Context, id uint) (*PollOption, error)
	DeleteOptions(ctx context.Context, pollID uint) error

	// 投票记录
	CreateVote(ctx context.Context, vote *PollVote) error
	CreateVotes(ctx context.Context, votes []*PollVote) (int, error) // 批量投票
	UpdateVote(ctx context.Context, vote *PollVote) error
	GetVoteByUnique(ctx context.Context, pollID, optionID, voterUserID uint) (*PollVote, error)
	ListVotesByPoll(ctx context.Context, pollID uint) ([]*PollVote, error)
	ListVotesByUserAndPoll(ctx context.Context, pollID, voterUserID uint) ([]*PollVote, error)
	CountVotersByPoll(ctx context.Context, pollID uint) (int64, error) // 投票人数

	// 结果汇总
	UpsertResult(ctx context.Context, result *PollResult) error
	GetResultsByPoll(ctx context.Context, pollID uint) ([]*PollResult, error)
	DeleteResults(ctx context.Context, pollID uint) error
}

type pollRepository struct {
	db *gorm.DB
}

func NewPollRepository(db *gorm.DB) PollRepository {
	return &pollRepository{db: db}
}

// ===== 投票 =====

func (r *pollRepository) CreatePoll(ctx context.Context, p *Poll) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pollRepository) GetPollByID(ctx context.Context, id uint) (*Poll, error) {
	var p Poll
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pollRepository) ListPolls(ctx context.Context, scopeType string, scopeID uint, page, pageSize int) ([]*Poll, int64, error) {
	var polls []*Poll
	var total int64
	query := r.db.WithContext(ctx).Model(&Poll{}).
		Where("scope_type = ? AND scope_id = ?", scopeType, scopeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).Find(&polls).Error
	return polls, total, err
}

func (r *pollRepository) UpdatePoll(ctx context.Context, p *Poll) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// ===== 选项 =====

func (r *pollRepository) CreateOptions(ctx context.Context, opts []*PollOption) (int, error) {
	result := r.db.WithContext(ctx).Create(opts)
	return int(result.RowsAffected), result.Error
}

func (r *pollRepository) GetOptionsByPoll(ctx context.Context, pollID uint) ([]*PollOption, error) {
	var opts []*PollOption
	err := r.db.WithContext(ctx).
		Where("poll_id = ?", pollID).
		Order("sort_order ASC, id ASC").
		Find(&opts).Error
	return opts, err
}

func (r *pollRepository) GetOptionByID(ctx context.Context, id uint) (*PollOption, error) {
	var opt PollOption
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&opt).Error
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func (r *pollRepository) DeleteOptions(ctx context.Context, pollID uint) error {
	return r.db.WithContext(ctx).Where("poll_id = ?", pollID).Delete(&PollOption{}).Error
}

// ===== 投票记录 =====

func (r *pollRepository) CreateVote(ctx context.Context, vote *PollVote) error {
	return r.db.WithContext(ctx).Create(vote).Error
}

func (r *pollRepository) CreateVotes(ctx context.Context, votes []*PollVote) (int, error) {
	result := r.db.WithContext(ctx).Create(votes)
	return int(result.RowsAffected), result.Error
}

func (r *pollRepository) UpdateVote(ctx context.Context, vote *PollVote) error {
	return r.db.WithContext(ctx).Save(vote).Error
}

func (r *pollRepository) GetVoteByUnique(ctx context.Context, pollID, optionID, voterUserID uint) (*PollVote, error) {
	var v PollVote
	err := r.db.WithContext(ctx).
		Where("poll_id = ? AND option_id = ? AND voter_user_id = ?", pollID, optionID, voterUserID).
		First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *pollRepository) ListVotesByPoll(ctx context.Context, pollID uint) ([]*PollVote, error) {
	var votes []*PollVote
	err := r.db.WithContext(ctx).
		Where("poll_id = ?", pollID).
		Find(&votes).Error
	return votes, err
}

func (r *pollRepository) ListVotesByUserAndPoll(ctx context.Context, pollID, voterUserID uint) ([]*PollVote, error) {
	var votes []*PollVote
	err := r.db.WithContext(ctx).
		Where("poll_id = ? AND voter_user_id = ?", pollID, voterUserID).
		Find(&votes).Error
	return votes, err
}

func (r *pollRepository) CountVotersByPoll(ctx context.Context, pollID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PollVote{}).
		Where("poll_id = ?", pollID).
		Distinct("voter_user_id").Count(&count).Error
	return count, err
}

// ===== 结果汇总 =====

func (r *pollRepository) UpsertResult(ctx context.Context, result *PollResult) error {
	result.CalculatedAt = time.Now()
	return r.db.WithContext(ctx).Save(result).Error
}

func (r *pollRepository) GetResultsByPoll(ctx context.Context, pollID uint) ([]*PollResult, error) {
	var results []*PollResult
	err := r.db.WithContext(ctx).
		Where("poll_id = ?", pollID).
		Order("option_id ASC").
		Find(&results).Error
	return results, err
}

func (r *pollRepository) DeleteResults(ctx context.Context, pollID uint) error {
	return r.db.WithContext(ctx).Where("poll_id = ?", pollID).Delete(&PollResult{}).Error
}
