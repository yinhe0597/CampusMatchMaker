package classrepo

import "time"

// Class 班级模型
type Class struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SchoolID        uint      `json:"school_id" gorm:"column:school_id;not null"`
	Grade           string    `json:"grade" gorm:"column:grade;type:varchar(20);not null"`
	Department      *string   `json:"department" gorm:"column:department;type:varchar(50)"`
	Name            string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Code            *string   `json:"code" gorm:"column:code;type:varchar(50);uniqueIndex:uk_code"`
	CreatorUserID   uint      `json:"creator_user_id" gorm:"column:creator_user_id;not null"`
	InviteCode      *string   `json:"invite_code" gorm:"column:invite_code;type:varchar(10);uniqueIndex:uk_invite_code"`
	TimetableStatus int       `json:"timetable_status" gorm:"column:timetable_status;not null;default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Class) TableName() string { return "classes" }

// ClassMember 班级成员模型
type ClassMember struct {
	ID       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ClassID  uint      `json:"class_id" gorm:"column:class_id;not null;uniqueIndex:uk_class_user"`
	UserID   uint      `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:uk_class_user"`
	Role     string    `json:"role" gorm:"column:role;type:varchar(20);not null;default:'member'"`
	Status   int       `json:"status" gorm:"column:status;not null;default:1"`
	JoinedAt time.Time `json:"joined_at" gorm:"column:joined_at;autoCreateTime"`
}

func (ClassMember) TableName() string { return "class_members" }

// 角色常量
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// 成员状态常量
const (
	MemberActive = 1
	MemberLeft   = 0
)

// MemberInfo 成员信息（含用户昵称）
type MemberInfo struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}
