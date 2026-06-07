package user

import (
	"time"

	"gorm.io/gorm"
)

// School 学校模型
type School struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Code      string    `json:"code" gorm:"column:code;type:varchar(20);not null;uniqueIndex:uk_code"`
	Province  *string   `json:"province" gorm:"column:province;type:varchar(20)"`
	City      *string   `json:"city" gorm:"column:city;type:varchar(20)"`
	Status    int       `json:"status" gorm:"column:status;not null;default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (School) TableName() string { return "schools" }

// User 用户模型
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StudentID    string         `json:"-" gorm:"column:student_id;type:varchar(256);not null;uniqueIndex:uk_student_id"`
	Nickname     string         `json:"nickname" gorm:"column:nickname;type:varchar(50);not null"`
	AvatarURL    *string        `json:"avatar_url" gorm:"column:avatar_url;type:varchar(500)"`
	Phone        *string        `json:"-" gorm:"column:phone;type:varchar(128)"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;type:varchar(255);not null"`
	SchoolID     *uint          `json:"school_id" gorm:"column:school_id"`
	Status       int            `json:"status" gorm:"column:status;not null;default:1"`
	PrivacyLevel int            `json:"privacy_level" gorm:"column:privacy_level;not null;default:1"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index:idx_deleted_at"`

	// 关联
	School *School `json:"school,omitempty" gorm:"foreignKey:SchoolID"`
}

func (User) TableName() string { return "users" }

// StudentAuth 学号认证模型
type StudentAuth struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint       `json:"user_id" gorm:"column:user_id;not null;index:idx_user_id"`
	StudentID  string     `json:"-" gorm:"column:student_id;type:varchar(256);not null"`
	SchoolID   uint       `json:"school_id" gorm:"column:school_id;not null"`
	AuthMethod string     `json:"auth_method" gorm:"column:auth_method;type:varchar(20);not null"`
	AuthStatus int        `json:"auth_status" gorm:"column:auth_status;not null;default:0"`
	VerifiedAt *time.Time `json:"verified_at" gorm:"column:verified_at"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (StudentAuth) TableName() string { return "student_auth" }
