package timetable

import "time"

// ClassTimetable 班级公共课表条目
type ClassTimetable struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ClassID           uint       `json:"class_id" gorm:"column:class_id;not null;index:idx_class_day"`
	DayOfWeek         int        `json:"day_of_week" gorm:"column:day_of_week;not null"`
	PeriodStart       int        `json:"period_start" gorm:"column:period_start;not null"`
	PeriodEnd         int        `json:"period_end" gorm:"column:period_end;not null"`
	CourseName        string     `json:"course_name" gorm:"column:course_name;type:varchar(100);not null"`
	Teacher           *string    `json:"teacher" gorm:"column:teacher;type:varchar(50)"`
	Room              *string    `json:"room" gorm:"column:room;type:varchar(50)"`
	ContributorUserID *uint      `json:"contributor_user_id" gorm:"column:contributor_user_id"`
	Version           int        `json:"version" gorm:"column:version;not null;default:1"`
	Status            int        `json:"status" gorm:"column:status;not null;default:1"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (ClassTimetable) TableName() string { return "class_timetables" }

// 公共课表状态
const (
	CTStatusActive   = 1 // 有效
	CTStatusDeleted  = 0 // 删除
	CTStatusReplaced = 2 // 已纠错替换
)

// PersonalTimetable 个人课表条目
type PersonalTimetable struct {
	ID                  uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID              uint       `json:"user_id" gorm:"column:user_id;not null;index:idx_user_day"`
	ClassID             *uint      `json:"class_id" gorm:"column:class_id"`
	DayOfWeek           int        `json:"day_of_week" gorm:"column:day_of_week;not null"`
	PeriodStart         int        `json:"period_start" gorm:"column:period_start;not null"`
	PeriodEnd           int        `json:"period_end" gorm:"column:period_end;not null"`
	CourseName          string     `json:"course_name" gorm:"column:course_name;type:varchar(100);not null"`
	Source              string     `json:"source" gorm:"column:source;type:varchar(20);not null"`
	RefClassTimetableID *uint      `json:"ref_class_timetable_id" gorm:"column:ref_class_timetable_id"`
	IsOverridden        int        `json:"is_overridden" gorm:"column:is_overridden;not null;default:0"`
	CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt           *time.Time `json:"deleted_at" gorm:"column:deleted_at;index"`
}

func (PersonalTimetable) TableName() string { return "personal_timetables" }

// 个人课表来源
const (
	PTInherited = "inherited"
	PTPersonal  = "personal"
)

// TimetableCorrection 课表纠错记录
type TimetableCorrection struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ClassTimetableID      uint       `json:"class_timetable_id" gorm:"column:class_timetable_id;not null;index:idx_ct_id"`
	ReporterUserID        uint       `json:"reporter_user_id" gorm:"column:reporter_user_id;not null"`
	CorrectionType        string     `json:"correction_type" gorm:"column:correction_type;type:varchar(20);not null"`
	Description           *string    `json:"description" gorm:"column:description;type:text"`
	SuggestedCourseName   *string    `json:"suggested_course_name" gorm:"column:suggested_course_name;type:varchar(100)"`
	SuggestedPeriodStart  *int       `json:"suggested_period_start" gorm:"column:suggested_period_start"`
	SuggestedPeriodEnd    *int       `json:"suggested_period_end" gorm:"column:suggested_period_end"`
	Status                int        `json:"status" gorm:"column:status;not null;default:0"`
	ReviewedBy            *uint      `json:"reviewed_by" gorm:"column:reviewed_by"`
	ReviewedAt            *time.Time `json:"reviewed_at" gorm:"column:reviewed_at"`
	CreatedAt             time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	ResolvedAt            *time.Time `json:"resolved_at" gorm:"column:resolved_at"`
}

func (TimetableCorrection) TableName() string { return "timetable_corrections" }

// 纠错状态
const (
	TCStatusPending  = 0 // 待审核
	TCStatusApproved = 1 // 已采纳
	TCStatusRejected = 2 // 已驳回
)

// 纠错类型
const (
	TCTypeError   = "error"
	TCTypeMissing = "missing"
)

// PeriodToMinutes 中国大学标准节次→分钟映射
// 第1节 08:00-08:45 => 480-525
// 第2节 08:50-09:35 => 530-575
// 第3节 09:55-10:40 => 595-640
// 第4节 10:45-11:30 => 645-690
// 第5节 11:35-12:20 => 695-740
// 第6节 13:00-13:45 => 780-825
// 第7节 13:50-14:35 => 830-875
// 第8节 14:40-15:25 => 880-925
// 第9节 15:30-16:15 => 930-975
// 第10节 16:20-17:05 => 980-1025
// 第11节 17:10-17:55 => 1030-1075
// 第12节 18:30-19:15 => 1110-1155
var periodToStart = map[int]int{
	1: 480, 2: 530, 3: 595, 4: 645, 5: 695, 6: 780,
	7: 830, 8: 880, 9: 930, 10: 980, 11: 1030, 12: 1110,
}
var periodToEnd = map[int]int{
	1: 525, 2: 575, 3: 640, 4: 690, 5: 740, 6: 825,
	7: 875, 8: 925, 9: 975, 10: 1025, 11: 1075, 12: 1155,
}

// PeriodStartToMinutes 节次开始转分钟
func PeriodStartToMinutes(period int) int {
	if v, ok := periodToStart[period]; ok {
		return v
	}
	return 0
}

// PeriodEndToMinutes 节次结束转分钟
func PeriodEndToMinutes(period int) int {
	if v, ok := periodToEnd[period]; ok {
		return v
	}
	return 0
}

// CheckPeriodRange 检查节次范围是否有效
func CheckPeriodRange(start, end int) bool {
	return start >= 1 && start <= 12 && end >= 1 && end <= 12 && start <= end
}
