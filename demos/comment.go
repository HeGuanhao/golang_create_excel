// comment.go 渲染"批注" sheet：普通批注与富文本批注。
package demos

import (
	"github.com/xuri/excelize/v2"
)

func registerComment(f *excelize.File) error {
	const sheet = "批注"
	if err := newSheet(f, sheet, "带批注的单元格右上角有红色三角标记，鼠标悬停即可查看", 6); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "A", "C", 14); err != nil {
		return err
	}

	products := [][]interface{}{
		{"产品", "单价(元)", "库存"},
		{"机械键盘", 399, 120},
		{"无线鼠标", 129, 0},
		{"显示器", 1299, 45},
		{"扩展坞", 459, 8},
	}
	if err := setRows(f, sheet, 1, 4, products); err != nil {
		return err
	}
	headerStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      fill(colorHeader),
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, "A4", "C4", headerStyle); err != nil {
		return err
	}

	// 普通文本批注
	if err := f.AddComment(sheet, excelize.Comment{
		Cell:   "B5",
		Author: "产品部",
		Text:   "单价含 13% 增值税，大客户另有折扣。",
		Width:  220,
		Height: 80,
	}); err != nil {
		return err
	}
	// 富文本批注：多段不同格式
	if err := f.AddComment(sheet, excelize.Comment{
		Cell:   "C6",
		Author: "仓库",
		Paragraph: []excelize.RichTextRun{
			{Text: "库存告急：", Font: &excelize.Font{Bold: true, Color: colorRed}},
			{Text: "当前库存为 0，", Font: &excelize.Font{}},
			{Text: "请尽快补货！", Font: &excelize.Font{Bold: true, Underline: "single"}},
		},
		Width:  220,
		Height: 100,
	}); err != nil {
		return err
	}
	// 低库存提醒批注
	if err := f.AddComment(sheet, excelize.Comment{
		Cell:   "C8",
		Author: "仓库",
		Text:   "库存低于安全线（10 件），已触发补货流程。",
	}); err != nil {
		return err
	}
	return nil
}
