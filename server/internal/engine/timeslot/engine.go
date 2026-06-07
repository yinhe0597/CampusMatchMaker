// Package timeslot 共同空闲时段计算引擎
// 纯函数实现，无外部依赖，不 import 任何业务包。
// 相同输入永远得到相同输出，不关心调用者是谁。
package timeslot

import "sort"

// rawSlot 内部使用的中间数据结构
type rawSlot struct {
	day   int
	index int
	count int
}

// CalculateFreeSlots 计算一组用户的共同空闲时段。
// 输入为每个用户的已占用时段列表，输出为按参与率降序排列的空闲时段列表。
// 注意：输入数据中不包含任何个人信息（如课程名、教室），只包含时间数据。
func CalculateFreeSlots(schedules []UserSchedule, cfg EngineConfig) []FreeSlotResult {
	if len(schedules) == 0 {
		return nil
	}

	granularity := cfg.SlotGranularity
	if granularity <= 0 {
		granularity = 30
	}

	// 构建每天的时间槽矩阵：day(1-7) × slotIndex → 有空人数
	slotsPerDay := (cfg.DayEndMinutes - cfg.DayStartMinutes) / granularity
	type slotKey struct {
		day   int
		index int
	}
	// 统计每个时间槽有多少人有空
	slotFreeCount := make(map[slotKey]int)

	// 初始化：每个槽位从 0 开始
	for day := 1; day <= 7; day++ {
		for idx := 0; idx < slotsPerDay; idx++ {
			slotFreeCount[slotKey{day, idx}] = 0
		}
	}

	// 对每个用户，计算其空闲槽位
	for _, schedule := range schedules {
		// 合并该用户的占用时段
		merged := mergeOverlapping(schedule.Slots)

		for day := 1; day <= 7; day++ {
			for idx := 0; idx < slotsPerDay; idx++ {
				slotStart := cfg.DayStartMinutes + idx*granularity
				slotEnd := slotStart + granularity

				if !isOccupied(day, slotStart, slotEnd, merged) {
					slotFreeCount[slotKey{day, idx}]++
				}
			}
		}
	}

	// 收集结果：只保留至少有一人有空的槽位，并合并连续相同人数的槽位
	var rawSlots []rawSlot
	for day := 1; day <= 7; day++ {
		for idx := 0; idx < slotsPerDay; idx++ {
			count := slotFreeCount[slotKey{day, idx}]
			if count > 0 {
				rawSlots = append(rawSlots, rawSlot{day, idx, count})
			}
		}
	}

	// 合并连续且参与人数相同的槽位
	merged_results := mergeConsecutiveSlots(rawSlots, granularity, cfg.DayStartMinutes)

	// 按参与率降序排序
	sort.Slice(merged_results, func(i, j int) bool {
		if merged_results[i].Rate != merged_results[j].Rate {
			return merged_results[i].Rate > merged_results[j].Rate
		}
		return merged_results[i].DayOfWeek < merged_results[j].DayOfWeek
	})

	return merged_results
}

// mergeOverlapping 合并重叠的已占用时段
func mergeOverlapping(slots []OccupiedSlot) []OccupiedSlot {
	if len(slots) == 0 {
		return nil
	}

	// 按天分组，再按开始时间排序
	grouped := make(map[int][]OccupiedSlot)
	for _, s := range slots {
		grouped[s.DayOfWeek] = append(grouped[s.DayOfWeek], s)
	}

	var result []OccupiedSlot
	for day, daySlots := range grouped {
		sort.Slice(daySlots, func(i, j int) bool {
			return daySlots[i].StartMinutes < daySlots[j].StartMinutes
		})

		merged := daySlots[0]
		for i := 1; i < len(daySlots); i++ {
			cur := daySlots[i]
			if cur.StartMinutes <= merged.EndMinutes {
				// 重叠或相邻，合并
				if cur.EndMinutes > merged.EndMinutes {
					merged.EndMinutes = cur.EndMinutes
				}
			} else {
				result = append(result, OccupiedSlot{day, merged.StartMinutes, merged.EndMinutes})
				merged = cur
			}
		}
		result = append(result, OccupiedSlot{day, merged.StartMinutes, merged.EndMinutes})
	}

	return result
}

// isOccupied 检查某个时段是否被占用
func isOccupied(day, slotStart, slotEnd int, occupied []OccupiedSlot) bool {
	for _, o := range occupied {
		if o.DayOfWeek != day {
			continue
		}
		// 时段有重叠
		if slotStart < o.EndMinutes && slotEnd > o.StartMinutes {
			return true
		}
	}
	return false
}

// mergeConsecutiveSlots 合并连续且参与人数相同的槽位为一个时段
func mergeConsecutiveSlots(raw []rawSlot, granularity, dayStart int) []FreeSlotResult {
	if len(raw) == 0 {
		return nil
	}

	// 先按天+index排序
	sort.Slice(raw, func(i, j int) bool {
		if raw[i].day != raw[j].day {
			return raw[i].day < raw[j].day
		}
		return raw[i].index < raw[j].index
	})

	var results []FreeSlotResult
	current := raw[0]
	startIdx := current.index

	for i := 1; i < len(raw); i++ {
		r := raw[i]
		// 同一天、连续槽位、相同人数 → 合并
		if r.day == current.day && r.index == current.index+1 && r.count == current.count {
			current.index = r.index
		} else {
			// 输出当前合并段
			results = append(results, FreeSlotResult{
				DayOfWeek:      current.day,
				StartMinutes:   dayStart + startIdx*granularity,
				EndMinutes:     dayStart + (current.index+1)*granularity,
				AvailableCount: current.count,
				TotalCount:     0, // 调用方补充
				Rate:           0, // 调用方补充
			})
			current = r
			startIdx = r.index
		}
	}

	// 输出最后一段
	results = append(results, FreeSlotResult{
		DayOfWeek:      current.day,
		StartMinutes:   dayStart + startIdx*granularity,
		EndMinutes:     dayStart + (current.index+1)*granularity,
		AvailableCount: current.count,
		TotalCount:     0,
		Rate:           0,
	})

	return results
}

// FillTotalCount 为结果填充总人数和参与率（供 service 层调用）
func FillTotalCount(results []FreeSlotResult, totalUsers int) []FreeSlotResult {
	for i := range results {
		results[i].TotalCount = totalUsers
		if totalUsers > 0 {
			results[i].Rate = float64(results[i].AvailableCount) / float64(totalUsers)
		}
	}
	return results
}
