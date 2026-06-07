package classrepo

import (
	"context"

	"gorm.io/gorm"
)

// ClassRepository 班级仓储接口
type ClassRepository interface {
	CreateClass(ctx context.Context, c *Class) error
	GetClassByID(ctx context.Context, id uint) (*Class, error)
	GetClassByInviteCode(ctx context.Context, code string) (*Class, error)
	UpdateClass(ctx context.Context, c *Class) error

	CreateMember(ctx context.Context, m *ClassMember) error
	GetMember(ctx context.Context, classID, userID uint) (*ClassMember, error)
	GetMemberRole(ctx context.Context, classID, userID uint) (string, error)
	IsMember(ctx context.Context, classID, userID uint) (bool, error)
	ListMembers(ctx context.Context, classID uint, page, pageSize int) ([]MemberInfo, int64, error)
	UpdateMemberStatus(ctx context.Context, classID, userID uint, status int) error
	CountMembers(ctx context.Context, classID uint) (int64, error)
	GetScopeMemberIDs(ctx context.Context, classID uint) ([]uint, error)
	ListUserClasses(ctx context.Context, userID uint) ([]*Class, int64, error)
}

type classRepository struct {
	db *gorm.DB
}

// NewClassRepository 创建班级仓储
func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) CreateClass(ctx context.Context, c *Class) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *classRepository) GetClassByID(ctx context.Context, id uint) (*Class, error) {
	var c Class
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *classRepository) GetClassByInviteCode(ctx context.Context, code string) (*Class, error) {
	var c Class
	err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *classRepository) UpdateClass(ctx context.Context, c *Class) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *classRepository) CreateMember(ctx context.Context, m *ClassMember) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *classRepository) GetMember(ctx context.Context, classID, userID uint) (*ClassMember, error) {
	var m ClassMember
	err := r.db.WithContext(ctx).Where("class_id = ? AND user_id = ?", classID, userID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *classRepository) GetMemberRole(ctx context.Context, classID, userID uint) (string, error) {
	var m ClassMember
	err := r.db.WithContext(ctx).Where("class_id = ? AND user_id = ?", classID, userID).First(&m).Error
	if err != nil {
		return "", err
	}
	return m.Role, nil
}

func (r *classRepository) IsMember(ctx context.Context, classID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ClassMember{}).
		Where("class_id = ? AND user_id = ? AND status = ?", classID, userID, MemberActive).
		Count(&count).Error
	return count > 0, err
}

func (r *classRepository) ListMembers(ctx context.Context, classID uint, page, pageSize int) ([]MemberInfo, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&ClassMember{}).
		Where("class_id = ? AND status = ?", classID, MemberActive).
		Count(&total)

	var members []MemberInfo
	err := r.db.WithContext(ctx).Table("class_members cm").
		Select("cm.id, cm.user_id, u.nickname, COALESCE(u.avatar_url, '') as avatar_url, cm.role, cm.joined_at").
		Joins("LEFT JOIN users u ON u.id = cm.user_id").
		Where("cm.class_id = ? AND cm.status = ?", classID, MemberActive).
		Order("cm.role = 'owner' DESC, cm.joined_at ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&members).Error
	if err != nil {
		return nil, 0, err
	}
	if members == nil {
		members = []MemberInfo{}
	}
	return members, total, nil
}

func (r *classRepository) UpdateMemberStatus(ctx context.Context, classID, userID uint, status int) error {
	return r.db.WithContext(ctx).Model(&ClassMember{}).
		Where("class_id = ? AND user_id = ?", classID, userID).
		Update("status", status).Error
}

func (r *classRepository) CountMembers(ctx context.Context, classID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ClassMember{}).
		Where("class_id = ? AND status = ?", classID, MemberActive).
		Count(&count).Error
	return count, err
}

func (r *classRepository) ListUserClasses(ctx context.Context, userID uint) ([]*Class, int64, error) {
	var total int64
	r.db.WithContext(ctx).Table("class_members cm").
		Joins("JOIN classes c ON c.id = cm.class_id").
		Where("cm.user_id = ? AND cm.status = ?", userID, MemberActive).
		Count(&total)

	var classes []*Class
	err := r.db.WithContext(ctx).
		Table("classes c").
		Joins("JOIN class_members cm ON cm.class_id = c.id").
		Where("cm.user_id = ? AND cm.status = ?", userID, MemberActive).
		Order("c.created_at DESC").
		Find(&classes).Error
	if err != nil {
		return nil, 0, err
	}
	if classes == nil {
		classes = []*Class{}
	}
	return classes, total, nil
}

func (r *classRepository) GetScopeMemberIDs(ctx context.Context, classID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).Model(&ClassMember{}).
		Where("class_id = ? AND status = ?", classID, MemberActive).
		Pluck("user_id", &ids).Error
	return ids, err
}
