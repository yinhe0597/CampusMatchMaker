package user

import (
	"context"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id uint) (*User, error)
	GetUserByStudentID(ctx context.Context, encryptedStudentID string) (*User, error)
	GetUserWithSchool(ctx context.Context, id uint) (*User, *School, error)
	CreateStudentAuth(ctx context.Context, auth *StudentAuth) error
	GetSchoolByID(ctx context.Context, id uint) (*School, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByStudentID(ctx context.Context, encryptedStudentID string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("student_id = ?", encryptedStudentID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserWithSchool(ctx context.Context, id uint) (*User, *School, error) {
	var u User
	err := r.db.WithContext(ctx).Preload("School").Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, nil, err
	}
	return &u, u.School, nil
}

func (r *userRepository) CreateStudentAuth(ctx context.Context, auth *StudentAuth) error {
	return r.db.WithContext(ctx).Create(auth).Error
}

func (r *userRepository) GetSchoolByID(ctx context.Context, id uint) (*School, error) {
	var school School
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&school).Error
	if err != nil {
		return nil, err
	}
	return &school, nil
}
