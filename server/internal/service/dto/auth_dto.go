package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	StudentID string `json:"student_id" binding:"required,min=6,max=30"`
	Password  string `json:"password" binding:"required,min=8,max=50"`
	Nickname  string `json:"nickname" binding:"required,min=1,max=50"`
	SchoolID  uint   `json:"school_id" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// RegisterResult 注册响应
type RegisterResult struct {
	UserID    uint   `json:"user_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// LoginResult 登录响应
type LoginResult struct {
	UserID    uint       `json:"user_id"`
	Nickname  string     `json:"nickname"`
	AvatarURL *string    `json:"avatar_url"`
	SchoolID  *uint      `json:"school_id"`
	Token     string     `json:"token"`
	ExpiresAt int64      `json:"expires_at"`
}

// TokenResult Token 刷新响应
type TokenResult struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// UserInfoResult 用户信息响应
type UserInfoResult struct {
	ID           uint    `json:"id"`
	Nickname     string  `json:"nickname"`
	StudentID    string  `json:"student_id"` // 脱敏后
	AvatarURL    *string `json:"avatar_url"`
	SchoolID     *uint   `json:"school_id"`
	SchoolName   *string `json:"school_name,omitempty"`
	Status       int     `json:"status"`
	PrivacyLevel int     `json:"privacy_level"`
}
