package main

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

const sheet = "Sheet1"

// colName 将0-based列索引转为Excel列名（A, B, C...）
func colName(col int) string {
	name, _ := excelize.ColumnNumberToName(col + 1)
	return name
}

// cell 拼接单元格地址，col 0-based，row 1-based
func cell(col, row int) string {
	return fmt.Sprintf("%s%d", colName(col), row)
}

func main() {
	now := time.Now()
	year, month := now.Year(), now.Month()
	today := now.Day()

	f := excelize.NewFile()
	defer f.Close()

	styles, err := NewStyles(f)
	if err != nil {
		fmt.Println("创建样式失败:", err)
		return
	}

	// 设置列宽（7列日历 + 左侧留白）
	for col := 0; col < 7; col++ {
		f.SetColWidth(sheet, colName(col), colName(col), 12)
	}

	// ── 第1行：年月标题（合并A1:G1）──
	titleCell := fmt.Sprintf("%s1:%s1", colName(0), colName(6))
	f.MergeCell(sheet, cell(0, 1), cell(6, 1))
	f.SetCellValue(sheet, cell(0, 1), fmt.Sprintf("%d年%02d月", year, int(month)))
	f.SetCellStyle(sheet, cell(0, 1), cell(6, 1), styles.Title)
	f.SetRowHeight(sheet, 1, 36)
	_ = titleCell

	// ── 第2行：星期表头 ──
	headers := WeekdayNames()
	for col, name := range headers {
		f.SetCellValue(sheet, cell(col, 2), name)
		f.SetCellStyle(sheet, cell(col, 2), cell(col, 2), styles.Header)
	}
	f.SetRowHeight(sheet, 2, 28)

	// ── 第3行起：日历数据 ──
	weeks := BuildCalendar(year, month)
	for wi, week := range weeks {
		row := wi + 3
		f.SetRowHeight(sheet, row, 32)
		for col, day := range week {
			c := cell(col, row)
			if day == 0 {
				f.SetCellValue(sheet, c, "")
				f.SetCellStyle(sheet, c, c, styles.Empty)
				continue
			}
			isWeekend := col >= 5 // 周六=5, 周日=6
			switch {
			case day == today:
				f.SetCellValue(sheet, c, day)
				f.SetCellStyle(sheet, c, c, styles.Today)
			case isWeekend:
				f.SetCellValue(sheet, c, day)
				f.SetCellStyle(sheet, c, c, styles.Weekend)
			default:
				f.SetCellValue(sheet, c, day)
				f.SetCellStyle(sheet, c, c, styles.Weekday)
			}
		}
	}

	// ── 信息栏：日历下方两行 ──
	infoRow := len(weeks) + 4
	f.SetRowHeight(sheet, infoRow, 24)
	f.SetRowHeight(sheet, infoRow+1, 24)

	// 生成时间
	f.MergeCell(sheet, cell(0, infoRow), cell(1, infoRow))
	f.SetCellValue(sheet, cell(0, infoRow), "生成时间")
	f.SetCellStyle(sheet, cell(0, infoRow), cell(1, infoRow), styles.InfoKey)
	f.MergeCell(sheet, cell(2, infoRow), cell(6, infoRow))
	f.SetCellValue(sheet, cell(2, infoRow), now.Format("2006-01-02 15:04:05"))
	f.SetCellStyle(sheet, cell(2, infoRow), cell(6, infoRow), styles.InfoVal)

	// 本月天数
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	f.MergeCell(sheet, cell(0, infoRow+1), cell(1, infoRow+1))
	f.SetCellValue(sheet, cell(0, infoRow+1), "本月天数")
	f.SetCellStyle(sheet, cell(0, infoRow+1), cell(1, infoRow+1), styles.InfoKey)
	f.MergeCell(sheet, cell(2, infoRow+1), cell(6, infoRow+1))
	f.SetCellValue(sheet, cell(2, infoRow+1), fmt.Sprintf("%d 天", daysInMonth))
	f.SetCellStyle(sheet, cell(2, infoRow+1), cell(6, infoRow+1), styles.InfoVal)

	fileName := "current_time.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println("保存失败:", err)
		return
	}
	fmt.Println("日历已生成:", fileName)
}
