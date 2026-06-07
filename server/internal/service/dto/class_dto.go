package dto

// CreateClassRequest 创建班级请求
type CreateClassRequest struct {
	SchoolID   uint    `json:"school_id" binding:"required"`
	Grade      string  `json:"grade" binding:"required"`
	Department *string `json:"department"`
	Name       string  `json:"name" binding:"required"`
}

// JoinClassRequest 加入班级请求
type JoinClassRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

// ClassDetailResult 班级详情响应
type ClassDetailResult struct {
	ID              uint           `json:"id"`
	SchoolID        uint           `json:"school_id"`
	Grade           string         `json:"grade"`
	Department      *string        `json:"department"`
	Name            string         `json:"name"`
	Code            *string        `json:"code"`
	CreatorUserID   uint           `json:"creator_user_id"`
	InviteCode      *string        `json:"invite_code"`
	TimetableStatus int            `json:"timetable_status"`
	MemberCount     int64          `json:"member_count"`
	MyRole          string         `json:"my_role"`
	CreatedAt       string         `json:"created_at"`
}

// CreateClassResult 创建班级响应
type CreateClassResult struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	InviteCode      string  `json:"invite_code"`
	TimetableStatus int     `json:"timetable_status"`
}

// JoinClassResult 加入班级响应
type JoinClassResult struct {
	ClassID   uint   `json:"class_id"`
	ClassName string `json:"class_name"`
	InviteCode string `json:"invite_code"`
	Role      string `json:"role"`
}

// MemberListResult 成员列表响应
type MemberListResult struct {
	Members []MemberItem `json:"list"`
	Total   int64        `json:"total"`
	Page    int          `json:"page"`
	PageSize int         `json:"page_size"`
}

type MemberItem struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

// MyClassItem 我的班级列表项
type MyClassItem struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Grade           string `json:"grade"`
	Department      string `json:"department"`
	InviteCode      string `json:"invite_code"`
	Role            string `json:"role"`
	MemberCount     int64  `json:"member_count"`
	TimetableStatus int    `json:"timetable_status"`
	CreatedAt       string `json:"created_at"`
}

// LookupClassResult 通过邀请码查找班级结果
type LookupClassResult struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Grade       string `json:"grade"`
	MemberCount int64  `json:"member_count"`
}
