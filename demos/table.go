// table.go 渲染"表格" sheet：AddTable 表格样式、冻结窗格、自动筛选。
package demos

import (
	"github.com/xuri/excelize/v2"
)

func registerTable(f *excelize.File) error {
	const sheet = "表格"
	if err := newSheet(f, sheet, "左侧为 AddTable 表格（首行已冻结），右侧为普通区域的自动筛选", 9); err != nil {
		return err
	}
	if err := f.SetColWidth(sheet, "A", "I", 12); err != nil {
		return err
	}

	// ── 员工表：AddTable ──
	employees := [][]interface{}{
		{"工号", "姓名", "部门", "职位", "月薪"},
		{"E001", "张伟", "研发部", "高级工程师", 26000},
		{"E002", "李娜", "市场部", "市场专员", 12000},
		{"E003", "王强", "研发部", "架构师", 35000},
		{"E004", "刘洋", "财务部", "会计", 11000},
		{"E005", "陈静", "人事部", "HR 经理", 18000},
		{"E006", "杨帆", "研发部", "测试工程师", 15000},
		{"E007", "赵敏", "市场部", "品牌经理", 20000},
		{"E008", "周杰", "财务部", "财务总监", 32000},
	}
	if err := setRows(f, sheet, 1, 4, employees); err != nil {
		return err
	}
	if err := f.AddTable(sheet, &excelize.Table{
		Range:     "A4:E12",
		Name:      "EmployeeTable",
		StyleName: "TableStyleMedium9",
	}); err != nil {
		return err
	}

	// 冻结前 4 行（标题 2 行 + 空行 + 表头），滚动时表头保持可见
	if err := f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      4,
		TopLeftCell: "A5",
		ActivePane:  "bottomLeft",
		Selection: []excelize.Selection{
			{SQRef: "A5", ActiveCell: "A5", Pane: "bottomLeft"},
		},
	}); err != nil {
		return err
	}

	// ── 右侧普通区域：AutoFilter 自动筛选 ──
	cities := [][]interface{}{
		{"城市", "人口(万)", "排名"},
		{"上海", 2487, 1},
		{"北京", 2189, 2},
		{"深圳", 1756, 3},
		{"广州", 1868, 4},
		{"成都", 2094, 5},
		{"杭州", 1220, 6},
	}
	if err := setRows(f, sheet, 7, 4, cities); err != nil {
		return err
	}
	cityHeader, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"},
		Fill:      fill(colorOrange),
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, "G4", "I4", cityHeader); err != nil {
		return err
	}
	if err := f.AutoFilter(sheet, "G4:I10", []excelize.AutoFilterOptions{}); err != nil {
		return err
	}
	return nil
}
