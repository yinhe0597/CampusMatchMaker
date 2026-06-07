// 中国大学标准节次时间映射
// 第1节 08:00-08:45
// 第2节 08:50-09:35
// 第3节 09:55-10:40
// 第4节 10:45-11:30
// 第5节 11:35-12:20
// 第6节 13:00-13:45
// 第7节 13:50-14:35
// 第8节 14:40-15:25
// 第9节 15:30-16:15
// 第10节 16:20-17:05
// 第11节 17:10-17:55
// 第12节 18:30-19:15

export const PERIOD_TIMES = [
  { period: 1, start: '08:00', end: '08:45' },
  { period: 2, start: '08:50', end: '09:35' },
  { period: 3, start: '09:55', end: '10:40' },
  { period: 4, start: '10:45', end: '11:30' },
  { period: 5, start: '11:35', end: '12:20' },
  { period: 6, start: '13:00', end: '13:45' },
  { period: 7, start: '13:50', end: '14:35' },
  { period: 8, start: '14:40', end: '15:25' },
  { period: 9, start: '15:30', end: '16:15' },
  { period: 10, start: '16:20', end: '17:05' },
  { period: 11, start: '17:10', end: '17:55' },
  { period: 12, start: '18:30', end: '19:15' },
]

export const DAY_LABELS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

export const DAY_LABELS_SHORT = ['一', '二', '三', '四', '五', '六', '日']

/**
 * 获取节次的时间显示文本
 */
export function getPeriodLabel(periodStart, periodEnd) {
  const s = PERIOD_TIMES.find(p => p.period === periodStart)
  const e = PERIOD_TIMES.find(p => p.period === periodEnd)
  if (!s || !e) return `第${periodStart}-${periodEnd}节`
  return `${s.start}-${e.end}`
}

/**
 * 获取节次范围跨度（占几行）
 */
export function getPeriodRowSpan(periodStart, periodEnd) {
  return periodEnd - periodStart + 1
}
