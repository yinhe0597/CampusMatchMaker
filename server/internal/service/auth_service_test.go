package service

import (
	"context"
	"testing"

	"campus_collab/internal/domain/user"
	"campus_collab/internal/infra/config"
	appErr "campus_collab/pkg/errors"
	"campus_collab/internal/service/dto"
	"campus_collab/pkg/utils"

	"gorm.io/gorm"
)

// ===== Mock UserRepository =====

type mockUserRepo struct {
	users  map[uint]*user.User
	schools map[uint]*user.School
	auths  []*user.StudentAuth
	nextID uint
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[uint]*user.User),
		schools: make(map[uint]*user.School),
		nextID:  1,
	}
}

func (m *mockUserRepo) CreateUser(ctx context.Context, u *user.User) error {
	u.ID = m.nextID
	m.nextID++
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetUserByStudentID(ctx context.Context, encryptedStudentID string) (*user.User, error) {
	for _, u := range m.users {
		if u.StudentID == encryptedStudentID {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) GetUserWithSchool(ctx context.Context, id uint) (*user.User, *user.School, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, nil, gorm.ErrRecordNotFound
	}
	var school *user.School
	if u.SchoolID != nil {
		s, ok := m.schools[*u.SchoolID]
		if ok {
			school = s
		}
	}
	return u, school, nil
}

func (m *mockUserRepo) CreateStudentAuth(ctx context.Context, auth *user.StudentAuth) error {
	m.auths = append(m.auths, auth)
	return nil
}

func (m *mockUserRepo) GetSchoolByID(ctx context.Context, id uint) (*user.School, error) {
	s, ok := m.schools[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return s, nil
}

// ===== Test Config =====

func testJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:       "test-secret-key-for-jwt-0123456789",
		ExpireHours:  1,
		RefreshHours: 7,
	}
}

// ===== Tests =====

func newTestAuthService() (*AuthService, *mockUserRepo) {
	repo := newMockUserRepo()
	repo.schools[1] = &user.School{ID: 1, Name: "测试大学"}
	svc := NewAuthService(repo, testJWTConfig(), "test-encrypt-key-0123456789")
	return svc, repo
}

// TestRegister_Success 注册成功
func TestRegister_Success(t *testing.T) {
	svc, _ := newTestAuthService()

	result, err := svc.Register(context.Background(), &dto.RegisterRequest{
		StudentID: "2024001",
		Password:  "password123",
		Nickname:  "小明",
		SchoolID:  1,
	})
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	if result.UserID == 0 {
		t.Error("UserID 不应为 0")
	}
	if result.Token == "" {
		t.Error("Token 不应为空")
	}
	if result.ExpiresAt == 0 {
		t.Error("ExpiresAt 不应为 0")
	}
}

// TestRegister_DuplicateStudentID 重复学号注册
func TestRegister_DuplicateStudentID(t *testing.T) {
	svc, repo := newTestAuthService()

	// 先注册一个用户
	hashedID := utils.HashStudentID("2024001")
	repo.users[1] = &user.User{ID: 1, StudentID: hashedID, Status: 1}

	_, err := svc.Register(context.Background(), &dto.RegisterRequest{
		StudentID: "2024001",
		Password:  "abc123",
		Nickname:  "小红",
		SchoolID:  1,
	})
	if err == nil {
		t.Error("重复学号应返回错误")
	}
}

// TestRegister_InvalidStudentID 无效学号格式
func TestRegister_InvalidStudentID(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.Register(context.Background(), &dto.RegisterRequest{
		StudentID: "ab", // 太短
		Password:  "test",
		Nickname:  "test",
		SchoolID:  1,
	})
	if err != appErr.ErrStudentIDFormat {
		t.Errorf("期望 ErrStudentIDFormat，实际: %v", err)
	}
}

// TestRegister_SchoolNotFound 学校不存在
func TestRegister_SchoolNotFound(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.Register(context.Background(), &dto.RegisterRequest{
		StudentID: "2024001",
		Password:  "test",
		Nickname:  "test",
		SchoolID:  999, // 不存在的学校
	})
	if err == nil {
		t.Error("学校不存在应返回错误")
	}
}

// TestLogin_Success 登录成功
func TestLogin_Success(t *testing.T) {
	svc, repo := newTestAuthService()

	// 创建测试用户
	hashedID := utils.HashStudentID("2024001")
	passwordHash, _ := utils.HashPassword("correct_password")
	repo.users[1] = &user.User{
		ID:           1,
		StudentID:    hashedID,
		Nickname:     "小明",
		PasswordHash: passwordHash,
		SchoolID:     func() *uint { v := uint(1); return &v }(),
		Status:       1,
	}

	result, err := svc.Login(context.Background(), &dto.LoginRequest{
		StudentID: "2024001",
		Password:  "correct_password",
	})
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	if result.UserID != 1 {
		t.Errorf("UserID 应为 1，实际: %d", result.UserID)
	}
	if result.Token == "" {
		t.Error("Token 不应为空")
	}
}

// TestLogin_WrongPassword 密码错误
func TestLogin_WrongPassword(t *testing.T) {
	svc, repo := newTestAuthService()

	hashedID := utils.HashStudentID("2024001")
	passwordHash, _ := utils.HashPassword("correct_password")
	repo.users[1] = &user.User{ID: 1, StudentID: hashedID, PasswordHash: passwordHash, Status: 1}

	_, err := svc.Login(context.Background(), &dto.LoginRequest{
		StudentID: "2024001",
		Password:  "wrong_password",
	})
	if err == nil {
		t.Error("密码错误应返回错误")
	}
}

// TestLogin_DisabledAccount 账号已禁用
func TestLogin_DisabledAccount(t *testing.T) {
	svc, repo := newTestAuthService()

	hashedID := utils.HashStudentID("2024001")
	passwordHash, _ := utils.HashPassword("password")
	repo.users[1] = &user.User{ID: 1, StudentID: hashedID, PasswordHash: passwordHash, Status: 0}

	_, err := svc.Login(context.Background(), &dto.LoginRequest{
		StudentID: "2024001",
		Password:  "password",
	})
	if err == nil {
		t.Error("已禁用账号应返回错误")
	}
}

// TestRefreshToken_Success 刷新Token成功
func TestRefreshToken_Success(t *testing.T) {
	svc, repo := newTestAuthService()
	repo.users[1] = &user.User{ID: 1, Nickname: "test"}

	result, err := svc.RefreshToken(context.Background(), 1)
	if err != nil {
		t.Fatalf("刷新Token失败: %v", err)
	}
	if result.Token == "" {
		t.Error("Token 不应为空")
	}
}

// TestRefreshToken_UserNotFound 用户不存在
func TestRefreshToken_UserNotFound(t *testing.T) {
	svc, _ := newTestAuthService()

	_, err := svc.RefreshToken(context.Background(), 999)
	if err == nil {
		t.Error("不存在的用户应返回错误")
	}
}

// TestGetCurrentUser_Success 获取当前用户
func TestGetCurrentUser_Success(t *testing.T) {
	svc, repo := newTestAuthService()

	repo.users[1] = &user.User{
		ID:           1,
		Nickname:     "小明",
		StudentID:    "abcdef1234567890", // 长哈希值
		SchoolID:     func() *uint { v := uint(1); return &v }(),
		Status:       1,
		PrivacyLevel: 1,
	}

	result, err := svc.GetCurrentUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCurrentUser 失败: %v", err)
	}
	if result.Nickname != "小明" {
		t.Errorf("Nickname 应为 '小明'，实际: %s", result.Nickname)
	}
	// 学号哈希应被脱敏
	if len(result.StudentID) >= 16 {
		t.Error("学号应被脱敏，不应显示完整哈希")
	}
	if result.SchoolName != nil && *result.SchoolName != "测试大学" {
		t.Errorf("学校名应为 '测试大学'，实际: %s", *result.SchoolName)
	}
}
