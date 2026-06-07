package dto

// ===== 公共课表 =====

type ClassTimetableEntry struct {
	DayOfWeek   int     `json:"day_of_week"`
	PeriodStart int     `json:"period_start"`
	PeriodEnd   int     `json:"period_end"`
	CourseName  string  `json:"course_name"`
	Teacher     *string `json:"teacher"`
	Room        *string `json:"room"`
}

type CreateClassTimetableRequest struct {
	Entries []ClassTimetableEntry `json:"entries" binding:"required,min=1"`
}

type CreateClassTimetableResult struct {
	CreatedCount    int  `json:"created_count"`
	TimetableStatus int  `json:"timetable_status"`
}

type ClassTimetableItem struct {
	ID              uint    `json:"id"`
	DayOfWeek       int     `json:"day_of_week"`
	PeriodStart     int     `json:"period_start"`
	PeriodEnd       int     `json:"period_end"`
	CourseName      string  `json:"course_name"`
	Teacher         *string `json:"teacher"`
	Room            *string `json:"room"`
	Version         int     `json:"version"`
}

type ClassTimetableResult struct {
	ClassID       uint                 `json:"class_id"`
	Entries       []ClassTimetableItem `json:"entries"`
	TotalEntries  int                  `json:"total_entries"`
}

type UpdateClassTimetableRequest struct {
	Entries []ClassTimetableEntry `json:"entries" binding:"required,min=1"`
}

type UpdateClassTimetableResult struct {
	UpdatedCount   int   `json:"updated_count"`
	AffectedMembers int64 `json:"affected_members"`
}

// ===== 个人课表 =====

type CreatePersonalTimetableRequest struct {
	ClassID     uint   `json:"class_id" binding:"required"`
	DayOfWeek   int    `json:"day_of_week" binding:"required,min=1,max=7"`
	PeriodStart int    `json:"period_start" binding:"required,min=1,max=12"`
	PeriodEnd   int    `json:"period_end" binding:"required,min=1,max=12"`
	CourseName  string `json:"course_name" binding:"required"`
}

type CreatePersonalTimetableResult struct {
	ID           uint   `json:"id"`
	Source       string `json:"source"`
	IsOverridden bool   `json:"is_overridden"`
}

type PersonalTimetableItem struct {
	ID           uint   `json:"id"`
	DayOfWeek    int    `json:"day_of_week"`
	PeriodStart  int    `json:"period_start"`
	PeriodEnd    int    `json:"period_end"`
	CourseName   string `json:"course_name"`
	Source       string `json:"source"`
	IsOverridden bool   `json:"is_overridden"`
}

type PersonalTimetableResult struct {
	UserID          uint                     `json:"user_id"`
	ClassID         uint                     `json:"class_id"`
	Entries         []PersonalTimetableItem  `json:"entries"`
	InheritedCount  int                      `json:"inherited_count"`
	PersonalCount   int                      `json:"personal_count"`
}

type UpdatePersonalTimetableRequest struct {
	CourseName   *string `json:"course_name"`
	IsOverridden *bool   `json:"is_overridden"`
}

// ===== 纠错 =====

type CreateCorrectionRequest struct {
	ClassTimetableID     uint    `json:"class_timetable_id" binding:"required"`
	CorrectionType       string  `json:"correction_type" binding:"required,oneof=error missing"`
	Description          *string `json:"description"`
	SuggestedCourseName  *string `json:"suggested_course_name"`
	SuggestedPeriodStart *int    `json:"suggested_period_start"`
	SuggestedPeriodEnd   *int    `json:"suggested_period_end"`
}

type CreateCorrectionResult struct {
	ID      uint   `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CorrectionItem struct {
	ID                   uint   `json:"id"`
	ClassTimetableID     uint   `json:"class_timetable_id"`
	ReporterUserID       uint   `json:"reporter_user_id"`
	ReporterNickname     string `json:"reporter_nickname"`
	CorrectionType       string `json:"correction_type"`
	Description          string `json:"description"`
	SuggestedCourseName  string `json:"suggested_course_name"`
	SuggestedPeriodStart *int   `json:"suggested_period_start"`
	SuggestedPeriodEnd   *int   `json:"suggested_period_end"`
	Status               int    `json:"status"`
	CreatedAt            string `json:"created_at"`
}

type CorrectionListResult struct {
	List     []CorrectionItem `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

type ReviewCorrectionRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
}
