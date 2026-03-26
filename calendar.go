package main

import "time"

// Week 代表日历中的一行（7天）
type Week [7]int // 0 表示非本月日期

// BuildCalendar 返回当月按周排列的日历，每行7天（周一到周日）
func BuildCalendar(year int, month time.Month) []Week {
	// 当月第一天
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	// 当月天数
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()

	// 周一为第0列，time.Weekday 中 Sunday=0，需要转换
	startCol := int(firstDay.Weekday()+6) % 7 // Mon=0 ... Sun=6

	var weeks []Week
	var week Week
	col := startCol
	for day := 1; day <= daysInMonth; day++ {
		week[col] = day
		col++
		if col == 7 {
			weeks = append(weeks, week)
			week = Week{}
			col = 0
		}
	}
	// 最后一行不满7天也加入
	if col > 0 {
		weeks = append(weeks, week)
	}
	return weeks
}

// WeekdayNames 返回表头（周一到周日）
func WeekdayNames() [7]string {
	return [7]string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
}
