package service

import (
	"context"
	"errors"
	"regexp"

	"campus_collab/internal/domain/user"
	"campus_collab/internal/infra/config"
	appErr "campus_collab/pkg/errors"
	"campus_collab/internal/service/dto"
	"campus_collab/pkg/utils"

	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	userRepo   user.UserRepository
	jwtCfg     config.JWTConfig
	encryptKey string
}

// NewAuthService 创建认证服务
func NewAuthService(repo user.UserRepository, jwtCfg config.JWTConfig, encryptKey string) *AuthService {
	return &AuthService{
		userRepo:   repo,
		jwtCfg:     jwtCfg,
		encryptKey: encryptKey,
	}
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResult, error) {
	// 1. 验证学号格式（字母+数字，6-30位）
	if !regexp.MustCompile(`^[a-zA-Z0-9]{6,30}$`).MatchString(req.StudentID) {
		return nil, appErr.ErrStudentIDFormat
	}

	// 2. 哈希学号（确定性，用于查找）
	hashedID := utils.HashStudentID(req.StudentID)

	// 3. 检查学号是否已注册
	existingUser, err := s.userRepo.GetUserByStudentID(ctx, hashedID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErr.Wrap(1000, "查询失败", err)
	}
	if existingUser != nil {
		return nil, appErr.ErrStudentIDExists
	}

	// 4. 验证学校存在
	_, err = s.userRepo.GetSchoolByID(ctx, req.SchoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.New(1002, "学校不存在")
		}
		return nil, appErr.Wrap(1000, "查询学校失败", err)
	}

	// 5. 哈希密码
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, appErr.Wrap(1000, "密码加密失败", err)
	}

	// 6. 创建用户
	newUser := &user.User{
		StudentID:    hashedID,
		Nickname:     req.Nickname,
		PasswordHash: hash,
		SchoolID:     &req.SchoolID,
		Status:       1,
		PrivacyLevel: 1,
	}
	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, appErr.Wrap(1000, "创建用户失败", err)
	}

	// 7. 创建认证记录
	auth := &user.StudentAuth{
		UserID:     newUser.ID,
		StudentID:  req.StudentID,
		SchoolID:   req.SchoolID,
		AuthMethod: "manual",
		AuthStatus: 0,
	}
	_ = s.userRepo.CreateStudentAuth(ctx, auth)

	// 8. 生成 JWT
	token, expiresAt, err := utils.GenerateToken(newUser.ID, req.StudentID, s.jwtCfg)
	if err != nil {
		return nil, appErr.Wrap(1000, "生成Token失败", err)
	}

	return &dto.RegisterResult{
		UserID:    newUser.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResult, error) {
	// 1. 哈希学号查找用户
	hashedID := utils.HashStudentID(req.StudentID)

	// 2. 查找用户
	u, err := s.userRepo.GetUserByStudentID(ctx, hashedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrPasswordWrong
		}
		return nil, appErr.Wrap(1000, "查询失败", err)
	}

	// 3. 验证密码
	if !utils.CheckPassword(req.Password, u.PasswordHash) {
		return nil, appErr.ErrPasswordWrong
	}

	// 4. 检查账号状态
	if u.Status != 1 {
		return nil, appErr.New(1107, "账号已被禁用")
	}

	// 5. 生成 JWT
	token, expiresAt, err := utils.GenerateToken(u.ID, req.StudentID, s.jwtCfg)
	if err != nil {
		return nil, appErr.Wrap(1000, "生成Token失败", err)
	}

	return &dto.LoginResult{
		UserID:    u.ID,
		Nickname:  u.Nickname,
		AvatarURL: u.AvatarURL,
		SchoolID:  u.SchoolID,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, userID uint) (*dto.TokenResult, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, appErr.ErrUnauthorized
	}

	// 解密学号用于生成新 token（token 中只存脱敏 ID）
	token, expiresAt, err := utils.GenerateToken(u.ID, "", s.jwtCfg)
	if err != nil {
		return nil, appErr.Wrap(1000, "生成Token失败", err)
	}

	return &dto.TokenResult{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*dto.UserInfoResult, error) {
	u, school, err := s.userRepo.GetUserWithSchool(ctx, userID)
	if err != nil {
		return nil, appErr.ErrUnauthorized
	}

	// 学号哈希脱敏：显示前6位
	maskedID := u.StudentID
	if len(maskedID) > 6 {
		maskedID = maskedID[:6] + "****"
	}

	result := &dto.UserInfoResult{
		ID:           u.ID,
		Nickname:     u.Nickname,
		StudentID:    maskedID,
		AvatarURL:    u.AvatarURL,
		SchoolID:     u.SchoolID,
		Status:       u.Status,
		PrivacyLevel: u.PrivacyLevel,
	}
	if school != nil {
		result.SchoolName = &school.Name
	}

	return result, nil
}

// maskStudentID 学号脱敏：前2后2，中间用 * 替代
func maskStudentID(id string) string {
	if len(id) <= 4 {
		return id
	}
	mask := ""
	for i := 2; i < len(id)-2; i++ {
		mask += "*"
	}
	return id[:2] + mask + id[len(id)-2:]
}
