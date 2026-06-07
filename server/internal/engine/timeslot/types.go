package timeslot

// OccupiedSlot 一个用户的单个已占用时段
// 不包含课程名、教室等业务信息，只有纯时间数据
type OccupiedSlot struct {
	DayOfWeek    int // 1=周一 ... 7=周日
	StartMinutes int // 从 00:00 起的分钟数，如 08:00 = 480
	EndMinutes   int // 从 00:00 起的分钟数，如 10:00 = 600
}

// UserSchedule 一个用户的完整日程
type UserSchedule struct {
	UserID string
	Slots  []OccupiedSlot
}

// FreeSlotResult 一个推荐的空闲时段
type FreeSlotResult struct {
	DayOfWeek      int     `json:"day_of_week"`      // 周几
	StartMinutes   int     `json:"start_minutes"`     // 开始时间（分钟）
	EndMinutes     int     `json:"end_minutes"`       // 结束时间（分钟）
	AvailableCount int     `json:"available_count"`   // 有空人数
	TotalCount     int     `json:"total_count"`       // 总人数
	Rate           float64 `json:"rate"`              // 参与率 (0.0 ~ 1.0)
}

// EngineConfig 引擎运行配置
type EngineConfig struct {
	DayStartMinutes int // 每天有效开始时间，默认 480 (08:00)
	DayEndMinutes   int // 每天有效结束时间，默认 1320 (22:00)
	SlotGranularity int // 时间粒度（分钟），默认 30
}

// DefaultConfig 返回默认引擎配置
func DefaultConfig() EngineConfig {
	return EngineConfig{
		DayStartMinutes: 480,  // 08:00
		DayEndMinutes:   1320, // 22:00
		SlotGranularity: 30,
	}
}
