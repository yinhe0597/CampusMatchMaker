package timetable

import (
	"context"

	"gorm.io/gorm"
)

// TimetableRepository 课表仓储接口
type TimetableRepository interface {
	// 班级公共课表
	CreateClassTimetables(ctx context.Context, entries []*ClassTimetable) (int, error)
	GetClassTimetables(ctx context.Context, classID uint) ([]*ClassTimetable, error)
	GetClassTimetableByID(ctx context.Context, id uint) (*ClassTimetable, error)
	DeleteClassTimetables(ctx context.Context, classID uint) error
	SoftDeleteClassTimetable(ctx context.Context, id uint) error
	UpdateClassTimetable(ctx context.Context, entry *ClassTimetable) error
	IncrementVersion(ctx context.Context, classID uint) error

	// 个人课表
	CreatePersonalTimetable(ctx context.Context, entry *PersonalTimetable) error
	CreatePersonalTimetables(ctx context.Context, entries []*PersonalTimetable) (int, error)
	GetPersonalTimetable(ctx context.Context, userID, classID uint) ([]*PersonalTimetable, error)
	GetPersonalTimetableByID(ctx context.Context, id uint) (*PersonalTimetable, error)
	UpdatePersonalTimetable(ctx context.Context, entry *PersonalTimetable) error
	SoftDeletePersonalTimetable(ctx context.Context, id uint) error
	DeletePersonalByClass(ctx context.Context, classID uint, source string) error
	CountPersonalByClass(ctx context.Context, classID uint) (int64, error)
	CheckTimeConflict(ctx context.Context, userID, classID uint, dayOfWeek, periodStart, periodEnd int, excludeID uint) (bool, error)

	// 继承
	HasInherited(ctx context.Context, userID, classID uint) (bool, error)

	// 纠错
	CreateCorrection(ctx context.Context, c *TimetableCorrection) error
	GetCorrectionByID(ctx context.Context, id uint) (*TimetableCorrection, error)
	ListCorrections(ctx context.Context, classID uint, status int, page, pageSize int) ([]*TimetableCorrection, int64, error)
	UpdateCorrection(ctx context.Context, c *TimetableCorrection) error
}

type timetableRepository struct {
	db *gorm.DB
}

func NewTimetableRepository(db *gorm.DB) TimetableRepository {
	return &timetableRepository{db: db}
}

// ===== 班级公共课表 =====

func (r *timetableRepository) CreateClassTimetables(ctx context.Context, entries []*ClassTimetable) (int, error) {
	result := r.db.WithContext(ctx).Create(entries)
	return int(result.RowsAffected), result.Error
}

func (r *timetableRepository) GetClassTimetables(ctx context.Context, classID uint) ([]*ClassTimetable, error) {
	var entries []*ClassTimetable
	err := r.db.WithContext(ctx).
		Where("class_id = ? AND status = ?", classID, CTStatusActive).
		Order("day_of_week ASC, period_start ASC").
		Find(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *timetableRepository) GetClassTimetableByID(ctx context.Context, id uint) (*ClassTimetable, error) {
	var entry ClassTimetable
	err := r.db.WithContext(ctx).Where("id = ? AND status = ?", id, CTStatusActive).First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *timetableRepository) DeleteClassTimetables(ctx context.Context, classID uint) error {
	return r.db.WithContext(ctx).Model(&ClassTimetable{}).
		Where("class_id = ?", classID).
		Update("status", CTStatusDeleted).Error
}

func (r *timetableRepository) SoftDeleteClassTimetable(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&ClassTimetable{}).
		Where("id = ?", id).
		Update("status", CTStatusDeleted).Error
}

func (r *timetableRepository) UpdateClassTimetable(ctx context.Context, entry *ClassTimetable) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *timetableRepository) IncrementVersion(ctx context.Context, classID uint) error {
	return r.db.WithContext(ctx).Model(&ClassTimetable{}).
		Where("class_id = ? AND status = ?", classID, CTStatusActive).
		Update("version", gorm.Expr("version + 1")).Error
}

// ===== 个人课表 =====

func (r *timetableRepository) CreatePersonalTimetable(ctx context.Context, entry *PersonalTimetable) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *timetableRepository) CreatePersonalTimetables(ctx context.Context, entries []*PersonalTimetable) (int, error) {
	result := r.db.WithContext(ctx).Create(entries)
	return int(result.RowsAffected), result.Error
}

func (r *timetableRepository) GetPersonalTimetable(ctx context.Context, userID, classID uint) ([]*PersonalTimetable, error) {
	var entries []*PersonalTimetable
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if classID > 0 {
		query = query.Where("class_id = ?", classID)
	}
	err := query.Order("day_of_week ASC, period_start ASC").Find(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *timetableRepository) GetPersonalTimetableByID(ctx context.Context, id uint) (*PersonalTimetable, error) {
	var entry PersonalTimetable
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *timetableRepository) UpdatePersonalTimetable(ctx context.Context, entry *PersonalTimetable) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *timetableRepository) SoftDeletePersonalTimetable(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PersonalTimetable{}, id).Error
}

func (r *timetableRepository) DeletePersonalByClass(ctx context.Context, classID uint, source string) error {
	query := r.db.WithContext(ctx).Where("class_id = ?", classID)
	if source != "" {
		query = query.Where("source = ?", source)
	}
	return query.Delete(&PersonalTimetable{}).Error
}

func (r *timetableRepository) CountPersonalByClass(ctx context.Context, classID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PersonalTimetable{}).
		Where("class_id = ? AND source = ?", classID, PTInherited).
		Count(&count).Error
	return count, err
}

func (r *timetableRepository) CheckTimeConflict(ctx context.Context, userID, classID uint, dayOfWeek, periodStart, periodEnd int, excludeID uint) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&PersonalTimetable{}).
		Where("user_id = ? AND day_of_week = ? AND period_start < ? AND period_end > ? AND deleted_at IS NULL",
			userID, dayOfWeek, periodEnd, periodStart)
	if classID > 0 {
		query = query.Where("class_id = ?", classID)
	}
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// ===== 继承 =====

func (r *timetableRepository) HasInherited(ctx context.Context, userID, classID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&PersonalTimetable{}).
		Where("user_id = ? AND class_id = ? AND source = ? AND deleted_at IS NULL",
			userID, classID, PTInherited).
		Count(&count).Error
	return count > 0, err
}

// ===== 纠错 =====

func (r *timetableRepository) CreateCorrection(ctx context.Context, c *TimetableCorrection) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *timetableRepository) GetCorrectionByID(ctx context.Context, id uint) (*TimetableCorrection, error) {
	var c TimetableCorrection
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *timetableRepository) ListCorrections(ctx context.Context, classID uint, status int, page, pageSize int) ([]*TimetableCorrection, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&TimetableCorrection{}).
		Joins("JOIN class_timetables ct ON ct.id = timetable_corrections.class_timetable_id")
	if classID > 0 {
		query = query.Where("ct.class_id = ?", classID)
	}
	if status >= 0 {
		query = query.Where("timetable_corrections.status = ?", status)
	}
	query.Count(&total)

	var corrections []*TimetableCorrection
	err := r.db.WithContext(ctx).
		Select("timetable_corrections.*").
		Joins("JOIN class_timetables ct ON ct.id = timetable_corrections.class_timetable_id").
		Where("ct.class_id = ?", classID).
		Order("timetable_corrections.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&corrections).Error
	if err != nil {
		return nil, 0, err
	}
	return corrections, total, nil
}

func (r *timetableRepository) UpdateCorrection(ctx context.Context, c *TimetableCorrection) error {
	return r.db.WithContext(ctx).Save(c).Error
}
