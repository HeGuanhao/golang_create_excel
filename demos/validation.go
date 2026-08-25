// validation.go 渲染"数据验证" sheet：下拉列表与数值范围校验。
package demos

import (
	"github.com/xuri/excelize/v2"
)

func registerValidation(f *excelize.File) error {
	const sheet = "数据验证"
	if err := newSheet(f, sheet, "B 列为下拉列表，D 列限制 18-65 的整数；选中单元格可查看输入提示", 6); err != nil {
		return err
	}
	for _, col := range []string{"A", "C"} {
		if err := f.SetColWidth(sheet, col, col, 18); err != nil {
			return err
		}
	}
	for _, col := range []string{"B", "D"} {
		if err := f.SetColWidth(sheet, col, col, 14); err != nil {
			return err
		}
	}

	headerStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      fill(colorHeader),
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}
	for _, hc := range []string{"A4", "B4", "C4", "D4"} {
		f.SetCellStyle(sheet, hc, hc, headerStyle)
	}
	f.SetCellValue(sheet, "A4", "部门（下拉列表）")
	f.SetCellValue(sheet, "B4", "选择")
	f.SetCellValue(sheet, "C4", "年龄（18-65 整数）")
	f.SetCellValue(sheet, "D4", "输入")

	// 预填几行示例值，展示合法输入的效果
	if err := setRows(f, sheet, 1, 5, [][]interface{}{
		{"研发部", "研发部", "年龄 28", 28},
		{"市场部", "市场部", "年龄 35", 35},
		{"财务部", nil, "年龄 41", 41},
	}); err != nil {
		return err
	}

	// 下拉列表验证：B5:B14
	dvList := excelize.NewDataValidation(true)
	dvList.SetSqref("B5:B14")
	if err := dvList.SetDropList([]string{"研发部", "市场部", "财务部", "人事部"}); err != nil {
		return err
	}
	dvList.SetInput("请选择部门", "从下拉列表中选择一个部门")
	dvList.SetError(excelize.DataValidationErrorStyleStop, "输入无效", "请从下拉列表中选择部门，不要手动输入。")
	if err := f.AddDataValidation(sheet, dvList); err != nil {
		return err
	}

	// 整数范围验证：D5:D14，18 到 65 之间
	dvAge := excelize.NewDataValidation(true)
	dvAge.SetSqref("D5:D14")
	if err := dvAge.SetRange(18, 65,
		excelize.DataValidationTypeWhole,
		excelize.DataValidationOperatorBetween); err != nil {
		return err
	}
	dvAge.SetInput("请输入年龄", "仅限 18 到 65 之间的整数")
	dvAge.SetError(excelize.DataValidationErrorStyleWarning, "年龄超出范围", "年龄须在 18 到 65 之间，请确认后继续。")
	if err := f.AddDataValidation(sheet, dvAge); err != nil {
		return err
	}
	return nil
}
