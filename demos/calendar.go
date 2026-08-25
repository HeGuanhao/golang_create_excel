// calendar.go 提供日历纯逻辑（可测试），并渲染"月历" sheet。
//
// 渲染沿用并改进了原有设计：年月标题栏、周一到周日表头、
// 周末红色高亮、今天橙色高亮、底部信息栏。
package demos

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// Week 代表日历中的一行（7 天），0 表示非本月日期。
type Week [7]int

// BuildCalendar 返回当月按周排列的日历，每行 7 天（周一到周日）。
func BuildCalendar(year int, month time.Month) []Week {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()

	// 周一为第 0 列；time.Weekday 中 Sunday=0，需要转换
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
	// 最后一行不满 7 天也加入
	if col > 0 {
		weeks = append(weeks, week)
	}
	return weeks
}

// WeekdayNames 返回表头（周一到周日）。
func WeekdayNames() [7]string {
	return [7]string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
}

// calendarStyles 是"月历" sheet 用到的全部样式。
type calendarStyles struct {
	Header  int // 星期表头（工作日）
	WkendHd int // 星期表头（周末，红色）
	Weekend int // 周末
	Weekday int // 工作日
	Today   int // 今天
	Empty   int // 空白格（非本月）
	InfoKey int // 信息栏 key
	InfoVal int // 信息栏 value
}

func newCalendarStyles(f *excelize.File) (*calendarStyles, error) {
	center := &excelize.Alignment{Horizontal: "center", Vertical: "center"}
	base := excelize.Style{Alignment: center, Border: borders(colorBorder, 1)}

	mk := func(font *excelize.Font, bg string, border []excelize.Border) (int, error) {
		s := base
		s.Font = font
		if bg != "" {
			s.Fill = fill(bg)
		}
		if border != nil {
			s.Border = border
		}
		return addStyle(f, &s)
	}

	var s calendarStyles
	var err error
	if s.Header, err = mk(&excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"}, colorHeader, borders("FFFFFF", 1)); err != nil {
		return nil, err
	}
	if s.WkendHd, err = mk(&excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"}, colorRed, borders("FFFFFF", 1)); err != nil {
		return nil, err
	}
	if s.Weekend, err = mk(&excelize.Font{Bold: true, Size: 12, Color: colorRed}, "FFE7E7", nil); err != nil {
		return nil, err
	}
	if s.Weekday, err = mk(&excelize.Font{Size: 12, Color: colorText}, colorLight, nil); err != nil {
		return nil, err
	}
	if s.Today, err = mk(&excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"}, colorOrange, borders("C55A11", 2)); err != nil {
		return nil, err
	}
	if s.Empty, err = mk(nil, "FFFFFF", borders("EEEEEE", 1)); err != nil {
		return nil, err
	}
	if s.InfoKey, err = mk(&excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"}, colorHeader, nil); err != nil {
		return nil, err
	}
	if s.InfoVal, err = mk(&excelize.Font{Size: 10, Color: colorText}, "EBF3FB", nil); err != nil {
		return nil, err
	}
	return &s, nil
}

// registerCalendar 渲染"月历" sheet：当月日历。
func registerCalendar(f *excelize.File) error {
	const sheet = "月历"
	now := time.Now()
	year, month, today := now.Year(), now.Month(), now.Day()

	if _, err := f.NewSheet(sheet); err != nil {
		return err
	}
	// 标题栏为"XXXX年XX月"，比默认 sheet 名更直观
	if err := writeTitle(f, sheet, 7,
		fmt.Sprintf("%d年%02d月", year, int(month)),
		"周一到周日排列；周末红色高亮，今天橙色高亮"); err != nil {
		return err
	}

	styles, err := newCalendarStyles(f)
	if err != nil {
		return err
	}

	// 列宽（7 列日历）
	for col := 1; col <= 7; col++ {
		name, _ := excelize.ColumnNumberToName(col)
		if err := f.SetColWidth(sheet, name, name, 12); err != nil {
			return err
		}
	}

	// 第 3 行：星期表头，周末两列用红色表头
	for i, name := range WeekdayNames() {
		c := cell(i+1, 3)
		f.SetCellValue(sheet, c, name)
		st := styles.Header
		if i >= 5 {
			st = styles.WkendHd
		}
		if err := f.SetCellStyle(sheet, c, c, st); err != nil {
			return err
		}
	}
	f.SetRowHeight(sheet, 3, 26)

	// 第 4 行起：日历数据
	weeks := BuildCalendar(year, month)
	for wi, week := range weeks {
		row := wi + 4
		f.SetRowHeight(sheet, row, 32)
		for col, day := range week {
			c := cell(col+1, row)
			st := styles.Empty
			if day != 0 {
				f.SetCellValue(sheet, c, day)
				switch {
				case day == today:
					st = styles.Today
				case col >= 5: // 周六、周日
					st = styles.Weekend
				default:
					st = styles.Weekday
				}
			}
			if err := f.SetCellStyle(sheet, c, c, st); err != nil {
				return err
			}
		}
	}

	// 日历下方：信息栏（生成时间、本月天数）
	infoRow := len(weeks) + 5
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	infos := [][2]string{
		{"生成时间", now.Format("2006-01-02 15:04:05")},
		{"本月天数", fmt.Sprintf("%d 天", daysInMonth)},
	}
	for i, kv := range infos {
		row := infoRow + i
		f.SetRowHeight(sheet, row, 22)
		if err := f.MergeCell(sheet, cell(1, row), cell(2, row)); err != nil {
			return err
		}
		f.SetCellValue(sheet, cell(1, row), kv[0])
		if err := f.SetCellStyle(sheet, cell(1, row), cell(2, row), styles.InfoKey); err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell(3, row), cell(7, row)); err != nil {
			return err
		}
		f.SetCellValue(sheet, cell(3, row), kv[1])
		if err := f.SetCellStyle(sheet, cell(3, row), cell(7, row), styles.InfoVal); err != nil {
			return err
		}
	}
	return nil
}
