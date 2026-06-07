package dto

import "time"

// ===== 创建投票 =====

type CreatePollRequest struct {
	Title         string       `json:"title" binding:"required"`
	Description   *string      `json:"description"`
	ScopeType     string       `json:"scope_type" binding:"required,oneof=class group"`
	ScopeID       uint         `json:"scope_id" binding:"required"`
	Deadline      *time.Time   `json:"deadline"`
	AutoRecommend bool         `json:"auto_recommend"`
	TimePreference  *TimePref  `json:"time_preference"`
}

type TimePref struct {
	DayStartHour      int `json:"day_start_hour"`      // 默认 8
	DayEndHour        int `json:"day_end_hour"`         // 默认 22
	MinDurationMin    int `json:"min_duration_minutes"` // 默认 60
	MaxRecommendations int `json:"max_recommendations"` // 默认 5
}

type CreatePollResult struct {
	PollID    uint   `json:"poll_id"`
	Status    string `json:"status"`
	OptionsCreated int `json:"options_created"`
}

// ===== 投票列表 =====

type PollListItem struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	CreatorUserID   uint       `json:"creator_user_id"`
	Status          string     `json:"status"`
	Deadline        *time.Time `json:"deadline"`
	MinParticipants int        `json:"min_participants"`
	VoterCount      int64      `json:"voter_count"`
	TotalOptions    int        `json:"total_options"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PollListResult struct {
	Polls      []*PollListItem `json:"polls"`
	TotalCount int64            `json:"total_count"`
}

// ===== 投票详情 =====

type PollOptionItem struct {
	ID                 uint     `json:"id"`
	SlotDate           string   `json:"slot_date"`
	SlotStartTime      string   `json:"slot_start_time"`
	SlotEndTime        string   `json:"slot_end_time"`
	DayOfWeek          int      `json:"day_of_week"`
	IsRecommended      int      `json:"is_recommended"`
	RecommendationRate *float64 `json:"recommendation_rate"`
	SortOrder          int      `json:"sort_order"`
}

type VoteRecord struct {
	OptionID uint   `json:"option_id"`
	Choice   string `json:"choice"`
}

type PollOptionResult struct {
	Option     PollOptionItem `json:"option"`
	YesCount   int            `json:"yes_count"`
	NoCount    int            `json:"no_count"`
	MaybeCount int            `json:"maybe_count"`
	TotalVotes int            `json:"total_votes"`
}

type PollDetailResult struct {
	ID              uint               `json:"id"`
	Title           string             `json:"title"`
	Description     *string            `json:"description"`
	ScopeType       string             `json:"scope_type"`
	ScopeID         uint               `json:"scope_id"`
	Status          string             `json:"status"`
	Deadline        *time.Time         `json:"deadline"`
	MinParticipants int                `json:"min_participants"`
	FinalOptionID   *uint              `json:"final_option_id"`
	CreatorUserID   uint               `json:"creator_user_id"`
	Options         []*PollOptionResult `json:"options"`
	MyVotes         []*VoteRecord      `json:"my_votes"`
	VoterCount      int64              `json:"voter_count"`
	CreatedAt       time.Time          `json:"created_at"`
}

// ===== 编辑投票 =====

type EditPollRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

// ===== 投票 =====

type SubmitVoteRequest struct {
	Votes []VoteInput `json:"votes" binding:"required,min=1"`
}

type VoteInput struct {
	OptionID uint   `json:"option_id" binding:"required"`
	Choice   string `json:"choice" binding:"required,oneof=yes no maybe"`
}

type SubmitVoteResult struct {
	VotedCount int `json:"voted_count"`
}

// ===== 投票结果 =====

type OptionResultItem struct {
	OptionID     uint    `json:"option_id"`
	SlotDate     string  `json:"slot_date"`
	SlotStartTime string `json:"slot_start_time"`
	SlotEndTime  string  `json:"slot_end_time"`
	DayOfWeek    int     `json:"day_of_week"`
	YesCount     int     `json:"yes_count"`
	NoCount      int     `json:"no_count"`
	MaybeCount   int     `json:"maybe_count"`
	TotalVotes   int     `json:"total_votes"`
}

type PollResultsResult struct {
	Items []*OptionResultItem `json:"items"`
}

// ===== 确认结果 =====

type FinalizePollRequest struct {
	FinalOptionID uint `json:"final_option_id" binding:"required"`
}
