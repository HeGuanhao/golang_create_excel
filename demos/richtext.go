// richtext.go 渲染"富文本" sheet：单元格内多格式混排。
package demos

import (
	"github.com/xuri/excelize/v2"
)

func registerRichText(f *excelize.File) error {
	const sheet = "富文本"
	if err := newSheet(f, sheet, "SetCellRichText 让同一个单元格内的文字拥有不同字体、颜色与上下标", 6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "A", "A", 16); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "B", "B", 56); err != nil {
		return err
	}

	labelStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: colorText},
		Fill:      fill(colorLight),
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    borders(colorBorder, 1),
	})
	if err != nil {
		return err
	}

	type richCase struct {
		label string
		runs  []excelize.RichTextRun
	}
	cases := []richCase{
		{"混排样式", []excelize.RichTextRun{
			{Text: "Excelize", Font: &excelize.Font{Bold: true, Size: 14, Color: colorBlue}},
			{Text: " 让 "},
			{Text: "Go", Font: &excelize.Font{Bold: true, Italic: true, Color: "00ADD8"}},
			{Text: " 操作 Excel 变得 "},
			{Text: "简单高效", Font: &excelize.Font{Bold: true, Color: colorOrange, Underline: "single"}},
			{Text: "。"},
		}},
		{"上下标", []excelize.RichTextRun{
			{Text: "水的化学式是 H"},
			{Text: "2", Font: &excelize.Font{VertAlign: "subscript"}},
			{Text: "O，质能方程是 E=mc"},
			{Text: "2", Font: &excelize.Font{VertAlign: "superscript"}},
			{Text: "。"},
		}},
		{"多色文本", []excelize.RichTextRun{
			{Text: "红色 ", Font: &excelize.Font{Color: colorRed, Bold: true}},
			{Text: "绿色 ", Font: &excelize.Font{Color: "63BE7B", Bold: true}},
			{Text: "蓝色 ", Font: &excelize.Font{Color: colorBlue, Bold: true}},
			{Text: "灰色", Font: &excelize.Font{Color: colorGray, Strike: true}},
		}},
	}

	row := 4
	for _, c := range cases {
		f.SetCellValue(sheet, cell(1, row), c.label)
		if err := f.SetCellStyle(sheet, cell(1, row), cell(1, row), labelStyle); err != nil {
			return err
		}
		if err := f.SetCellRichText(sheet, cell(2, row), c.runs); err != nil {
			return err
		}
		f.SetRowHeight(sheet, row, 26)
		row += 2 // 行间留白
	}
	return nil
}
