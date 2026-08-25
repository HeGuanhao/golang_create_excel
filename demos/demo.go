// Package demos 汇集 showcase.xlsx 中各特效 sheet 的渲染逻辑。
//
// 每个 demo 对应一个 sheet，通过 All 注册；demo 之间共享本文件中的
// 主题色与样式小工具，避免重复的样式代码。
package demos

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Demo 描述一个特效 sheet。
type Demo struct {
	Name     string                       // sheet 名称
	Desc     string                       // 演示内容简述
	Register func(f *excelize.File) error // 渲染函数
}

// All 返回全部 demo，顺序即生成工作簿中的 sheet 顺序。
func All() []Demo {
	return []Demo{
		{Name: "月历", Desc: "当月日历：标题、周末高亮、今天高亮", Register: registerCalendar},
		{Name: "样式一览", Desc: "字体 / 填充 / 边框 / 对齐 / 数字格式", Register: registerStyles},
		{Name: "条件格式", Desc: "数据条 / 色阶 / 图标集 / 高亮规则", Register: registerConditional},
		{Name: "图表", Desc: "柱状图 / 折线图 / 饼图", Register: registerChart},
		{Name: "表格", Desc: "AddTable / 冻结窗格 / 自动筛选", Register: registerTable},
		{Name: "富文本", Desc: "单元格内多格式混排", Register: registerRichText},
		{Name: "数据验证", Desc: "下拉列表与数值范围校验", Register: registerValidation},
		{Name: "图片", Desc: "内存中生成 PNG 并嵌入", Register: registerPicture},
		{Name: "批注", Desc: "普通批注与富文本批注", Register: registerComment},
	}
}

// 主题色板，各 demo 统一取用，保证视觉一致。
const (
	colorBlue   = "2E75B6" // 主色：标题栏
	colorHeader = "4472C4" // 表头蓝
	colorOrange = "ED7D31" // 强调橙（今天 / 重点）
	colorRed    = "C00000" // 周末红
	colorText   = "333333" // 正文深灰
	colorGray   = "808080" // 说明文字灰
	colorBorder = "D9D9D9" // 浅灰边框
	colorLight  = "F2F2F2" // 浅灰底纹
)

// fill 返回纯色填充。
func fill(color string) excelize.Fill {
	return excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1}
}

// borders 返回四边同色的边框，style 为 excelize 边框线型（1 细线、2 中线等）。
func borders(color string, style int) []excelize.Border {
	sides := []string{"left", "right", "top", "bottom"}
	bs := make([]excelize.Border, 0, len(sides))
	for _, s := range sides {
		bs = append(bs, excelize.Border{Type: s, Color: color, Style: style})
	}
	return bs
}

// addStyle 创建样式并包装错误信息。
func addStyle(f *excelize.File, s *excelize.Style) (int, error) {
	id, err := f.NewStyle(s)
	if err != nil {
		return 0, fmt.Errorf("创建样式失败: %w", err)
	}
	return id, nil
}

// cell 返回 1-based 行列对应的单元格地址，如 cell(1, 1) 为 "A1"。
func cell(col, row int) string {
	c, _ := excelize.CoordinatesToCellName(col, row)
	return c
}

// newSheet 新建 demo sheet 并写入统一标题栏（第 1 行）与说明栏（第 2 行）。
// lastCol 为标题合并到的最后一列（1-based），demo 内容建议从第 4 行开始。
func newSheet(f *excelize.File, name, desc string, lastCol int) error {
	if _, err := f.NewSheet(name); err != nil {
		return fmt.Errorf("新建工作表 %q 失败: %w", name, err)
	}
	return writeTitle(f, name, lastCol, name, desc)
}

// writeTitle 在 sheet 顶部写入统一的标题栏与说明栏。
func writeTitle(f *excelize.File, sheet string, lastCol int, title, desc string) error {
	titleStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Color: "FFFFFF"},
		Fill:      fill(colorBlue),
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}
	descStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Size: 10, Italic: true, Color: colorGray},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return err
	}

	tl, tr := cell(1, 1), cell(lastCol, 1)
	if err := f.MergeCell(sheet, tl, tr); err != nil {
		return err
	}
	f.SetCellValue(sheet, tl, title)
	if err := f.SetCellStyle(sheet, tl, tr, titleStyle); err != nil {
		return err
	}

	dl, dr := cell(1, 2), cell(lastCol, 2)
	if err := f.MergeCell(sheet, dl, dr); err != nil {
		return err
	}
	f.SetCellValue(sheet, dl, desc)
	if err := f.SetCellStyle(sheet, dl, dr, descStyle); err != nil {
		return err
	}

	f.SetRowHeight(sheet, 1, 32)
	f.SetRowHeight(sheet, 2, 20)
	return nil
}

// setRows 从 (col, row)（均 1-based）开始逐行写入数据。
func setRows(f *excelize.File, sheet string, col, row int, rows [][]interface{}) error {
	for i, r := range rows {
		if err := f.SetSheetRow(sheet, cell(col, row+i), &r); err != nil {
			return err
		}
	}
	return nil
}
