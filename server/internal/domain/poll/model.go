package poll

import "time"

// Poll 投票
type Poll struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatorUserID   uint       `json:"creator_user_id" gorm:"column:creator_user_id;not null;index:idx_creator"`
	Title           string     `json:"title" gorm:"column:title;type:varchar(200);not null"`
	Description     *string    `json:"description" gorm:"column:description;type:text"`
	ScopeType       string     `json:"scope_type" gorm:"column:scope_type;type:varchar(20);not null;index:idx_scope"`
	ScopeID         uint       `json:"scope_id" gorm:"column:scope_id;not null;index:idx_scope"`
	Status          string     `json:"status" gorm:"column:status;type:varchar(20);not null;default:'draft';index:idx_status"`
	Deadline        *time.Time `json:"deadline" gorm:"column:deadline"`
	MinParticipants int        `json:"min_participants" gorm:"column:min_participants;default:2"`
	FinalOptionID   *uint      `json:"final_option_id" gorm:"column:final_option_id"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	ClosedAt        *time.Time `json:"closed_at" gorm:"column:closed_at"`
}

func (Poll) TableName() string { return "polls" }

// 投票状态常量
const (
	PollStatusDraft     = "draft"
	PollStatusOpen      = "open"
	PollStatusClosed    = "closed"
	PollStatusFinalized = "finalized"
)

// PollOption 投票选项
type PollOption struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PollID             uint      `json:"poll_id" gorm:"column:poll_id;not null;index:idx_poll_sort"`
	SlotDate           string    `json:"slot_date" gorm:"column:slot_date;type:date;not null"`
	SlotStartTime      string    `json:"slot_start_time" gorm:"column:slot_start_time;type:time;not null"`
	SlotEndTime        string    `json:"slot_end_time" gorm:"column:slot_end_time;type:time;not null"`
	DayOfWeek          int       `json:"day_of_week" gorm:"column:day_of_week;default:0"` // 用于存储引擎推荐的原始星期
	IsRecommended      int       `json:"is_recommended" gorm:"column:is_recommended;not null;default:0"`
	RecommendationRate *float64  `json:"recommendation_rate" gorm:"column:recommendation_rate;type:decimal(5,4)"`
	SortOrder          int       `json:"sort_order" gorm:"column:sort_order;not null;default:0"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (PollOption) TableName() string { return "poll_options" }

// PollVote 投票记录
type PollVote struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PollID      uint      `json:"poll_id" gorm:"column:poll_id;not null;index:idx_poll_voter"`
	OptionID    uint      `json:"option_id" gorm:"column:option_id;not null"`
	VoterUserID uint      `json:"voter_user_id" gorm:"column:voter_user_id;not null;uniqueIndex:uk_poll_option_voter"`
	Choice      string    `json:"choice" gorm:"column:choice;type:varchar(10);not null"`
	VotedAt     time.Time `json:"voted_at" gorm:"column:voted_at;autoCreateTime"`
}

func (PollVote) TableName() string { return "poll_votes" }

// PollResult 投票结果汇总
type PollResult struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PollID            uint      `json:"poll_id" gorm:"column:poll_id;not null;uniqueIndex:uk_poll_option"`
	OptionID          uint      `json:"option_id" gorm:"column:option_id;not null;uniqueIndex:uk_poll_option"`
	YesCount          int       `json:"yes_count" gorm:"column:yes_count;not null;default:0"`
	NoCount           int       `json:"no_count" gorm:"column:no_count;not null;default:0"`
	MaybeCount        int       `json:"maybe_count" gorm:"column:maybe_count;not null;default:0"`
	TotalVotes        int       `json:"total_votes" gorm:"column:total_votes;not null;default:0"`
	ParticipationRate float64   `json:"participation_rate" gorm:"column:participation_rate;type:decimal(5,4);not null;default:0"`
	CalculatedAt      time.Time `json:"calculated_at" gorm:"column:calculated_at;not null"`
}

func (PollResult) TableName() string { return "poll_results" }

// 投票选项
const (
	ChoiceYes   = "yes"
	ChoiceNo    = "no"
	ChoiceMaybe = "maybe"
)
