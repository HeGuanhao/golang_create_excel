// calendar.go 保留原有的日历纯逻辑 API。
//
// 说明：Go 不允许子包导入 main 包，因此日历逻辑的实际实现位于
// demos/calendar.go（同时被"月历" sheet 渲染复用），这里以类型别名和
// 转发函数的形式保留原有 API，便于根包测试直接调用。
package main

import (
	"time"

	"hgh/golang_create_excel/demos"
)

// Week 代表日历中的一行（7 天），0 表示非本月日期。
type Week = demos.Week

// BuildCalendar 返回当月按周排列的日历，每行 7 天（周一到周日）。
func BuildCalendar(year int, month time.Month) []Week {
	return demos.BuildCalendar(year, month)
}

// WeekdayNames 返回表头（周一到周日）。
func WeekdayNames() [7]string {
	return demos.WeekdayNames()
}
