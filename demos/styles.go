// styles.go 渲染"样式一览" sheet：字体、填充、边框、对齐、数字格式。
package demos

import (
	"github.com/xuri/excelize/v2"
)

// styleCase 描述一行样式示例：A 列为说明，B 列为应用了样式的示例值。
type styleCase struct {
	label string
	value interface{}
	style *excelize.Style
}

func registerStyles(f *excelize.File) error {
	const sheet = "样式一览"
	if err := newSheet(f, sheet, "每一行演示一种单元格样式，A 列为说明、B 列为效果", 4); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "A", "A", 22); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "B", "B", 34); err != nil {
		return err
	}

	center := &excelize.Alignment{Horizontal: "center", Vertical: "center"}
	labelStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: colorText},
		Fill:      fill(colorLight),
		Alignment: center,
		Border:    borders(colorBorder, 1),
	})
	if err != nil {
		return err
	}

	// 分组标题 -> 该组示例
	numFmt := func(code string) *excelize.Style {
		return &excelize.Style{CustomNumFmt: &code, Alignment: center, Border: borders(colorBorder, 1)}
	}
	groups := []struct {
		title string
		cases []styleCase
	}{
		{"字体", []styleCase{
			{"加粗", "加粗文字", &excelize.Style{Font: &excelize.Font{Bold: true}}},
			{"斜体", "斜体文字", &excelize.Style{Font: &excelize.Font{Italic: true}}},
			{"下划线", "下划线文字", &excelize.Style{Font: &excelize.Font{Underline: "single"}}},
			{"删除线", "删除线文字", &excelize.Style{Font: &excelize.Font{Strike: true}}},
			{"彩色 + 字号", "16 号红色字", &excelize.Style{Font: &excelize.Font{Size: 16, Color: colorRed}}},
			{"指定字体", "微软雅黑", &excelize.Style{Font: &excelize.Font{Family: "微软雅黑", Size: 12}}},
		}},
		{"填充", []styleCase{
			{"蓝色填充白字", "深蓝", &excelize.Style{Font: &excelize.Font{Color: "FFFFFF", Bold: true}, Fill: fill(colorBlue), Alignment: center}},
			{"橙色填充白字", "橙色", &excelize.Style{Font: &excelize.Font{Color: "FFFFFF", Bold: true}, Fill: fill(colorOrange), Alignment: center}},
			{"浅灰填充", "浅灰", &excelize.Style{Fill: fill(colorLight), Alignment: center}},
		}},
		{"边框", []styleCase{
			{"细边框", "细线（Style 1）", &excelize.Style{Border: borders(colorText, 1), Alignment: center}},
			{"中边框", "中线（Style 2）", &excelize.Style{Border: borders(colorBlue, 2), Alignment: center}},
			{"虚线边框", "虚线（Style 3）", &excelize.Style{Border: borders(colorRed, 3), Alignment: center}},
		}},
		{"对齐", []styleCase{
			{"左对齐", "靠左", &excelize.Style{Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"}}},
			{"居中", "居中", &excelize.Style{Alignment: center}},
			{"右对齐", "靠右", &excelize.Style{Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"}}},
			{"自动换行", "这是一段很长很长的说明文字，开启自动换行后会在单元格内折行显示。", &excelize.Style{Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"}}},
			{"旋转 45 度", "斜着看", &excelize.Style{Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", TextRotation: 45}}},
		}},
		{"数字格式", []styleCase{
			{"千分位", 1234567.891, numFmt("#,##0.00")},
			{"百分比", 0.856, numFmt("0.0%")},
			{"货币", 9988.5, numFmt("¥#,##0.00")},
			{"日期", 45292, numFmt("yyyy\"年\"mm\"月\"dd\"日\"")},
			{"科学计数", 123456789, numFmt("0.00E+00")},
		}},
	}

	row := 4
	for _, g := range groups {
		// 分组标题：合并 A:B，蓝底白字
		secStyle, err := addStyle(f, &excelize.Style{
			Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
			Fill:      fill(colorHeader),
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		})
		if err != nil {
			return err
		}
		if err := f.MergeCell(sheet, cell(1, row), cell(2, row)); err != nil {
			return err
		}
		f.SetCellValue(sheet, cell(1, row), "■ "+g.title)
		if err := f.SetCellStyle(sheet, cell(1, row), cell(2, row), secStyle); err != nil {
			return err
		}
		f.SetRowHeight(sheet, row, 22)
		row++

		for _, c := range g.cases {
			st, err := addStyle(f, c.style)
			if err != nil {
				return err
			}
			f.SetCellValue(sheet, cell(1, row), c.label)
			if err := f.SetCellStyle(sheet, cell(1, row), cell(1, row), labelStyle); err != nil {
				return err
			}
			f.SetCellValue(sheet, cell(2, row), c.value)
			if err := f.SetCellStyle(sheet, cell(2, row), cell(2, row), st); err != nil {
				return err
			}
			f.SetRowHeight(sheet, row, 22)
			row++
		}
		row++ // 组间空一行
	}
	return nil
}
