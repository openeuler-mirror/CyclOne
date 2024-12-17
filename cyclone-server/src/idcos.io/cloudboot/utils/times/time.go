package times

import "time"

// LatestWeek 返回最近一周的时间列表
func LatestWeek(base time.Time, asc bool) []time.Time {
	return LatestDays(base, 7, asc)
}

// LatestMonth 返回最近一个月的时间列表
func LatestMonth(base time.Time, asc bool) []time.Time {
	return LatestDays(base, 30, asc)
}

// LatestDays 返回最近N天的时间列表
func LatestDays(base time.Time, days int, asc bool) []time.Time {
	if base.IsZero() || days <= 0 {
		return nil
	}
	items := make([]time.Time, 0, days)
	if asc {
		// 日期递增
		base = base.Add(time.Hour * -24 * time.Duration(days-1))
		for i := 0; i < days; i++ {
			items = append(items, base)
			base = base.Add(time.Hour * 24)
		}
	} else {
		// 日期递减
		for i := 0; i < days; i++ {
			items = append(items, base)
			base = base.Add(time.Hour * -24)
		}
	}

	return items
}
