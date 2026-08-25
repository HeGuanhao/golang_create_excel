// conditional.go 渲染"条件格式" sheet：数据条、色阶、图标集、高亮规则。
package demos

import (
	"github.com/xuri/excelize/v2"
)

func registerConditional(f *excelize.File) error {
	const sheet = "条件格式"
	if err := newSheet(f, sheet, "四组数据分别应用数据条、三阶色阶、交通灯图标集与数值高亮", 7); err != nil {
		return err
	}
	for _, col := range []string{"A", "C", "E", "G"} {
		if err := f.SetColWidth(sheet, col, col, 16); err != nil {
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

	headerRow, dataStart, dataEnd := 4, 5, 15

	// 高亮规则需要的格式（>90 分时绿底白字加粗）
	highlight, err := addStyle(f, &excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: fill("63BE7B"),
	})
	if err != nil {
		return err
	}

	// 四组演示数据：表头 + 数值 + 条件格式规则
	type block struct {
		col    int
		header string
		values []interface{}
		rules  []excelize.ConditionalFormatOptions
	}
	blks := []block{
		{1, "成绩 · 数据条", []interface{}{95, 82, 67, 73, 58, 90, 76, 64, 88, 49, 71},
			[]excelize.ConditionalFormatOptions{{
				Type:           "data_bar",
				Criteria:       "=", // data_bar 无判据，但 v2.8.1 要求给出合法 Criteria 占位
				MinType:        "min",
				MaxType:        "max",
				BarColor:       colorBlue,
				BarBorderColor: colorBlue,
				BarSolid:       true,
			}}},
		{3, "温度 · 三阶色阶", []interface{}{-5, 3, 12, 18, 24, 31, 36, 28, 19, 9, 1},
			[]excelize.ConditionalFormatOptions{{
				Type:     "3_color_scale",
				Criteria: "=", // 同上，占位
				MinType:  "min", MinColor: "63BE7B",
				MidType: "percentile", MidValue: "50", MidColor: "FFEB84",
				MaxType: "max", MaxColor: "F8696B",
			}}},
		{5, "完成率 · 图标集", []interface{}{0.95, 0.62, 0.30, 0.78, 0.45, 0.88, 0.15, 0.55, 0.70, 0.25, 0.99},
			[]excelize.ConditionalFormatOptions{{
				Type:      "icon_set",
				IconStyle: "3TrafficLights1",
			}}},
		{7, "分数 · 高亮 >90", []interface{}{88, 91, 76, 95, 60, 99, 83, 92, 71, 85, 97},
			[]excelize.ConditionalFormatOptions{{
				Type:     "cell",
				Criteria: ">",
				Value:    "90",
				Format:   highlight,
			}}},
	}

	centerStyle, err := addStyle(f, &excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}

	for _, b := range blks {
		hc := cell(b.col, headerRow)
		f.SetCellValue(sheet, hc, b.header)
		if err := f.SetCellStyle(sheet, hc, hc, headerStyle); err != nil {
			return err
		}
		for i, v := range b.values {
			c := cell(b.col, dataStart+i)
			f.SetCellValue(sheet, c, v)
			if err := f.SetCellStyle(sheet, c, c, centerStyle); err != nil {
				return err
			}
		}
		// 完成率一列按百分比显示
		if b.col == 5 {
			pct, err := addStyle(f, &excelize.Style{
				CustomNumFmt: strPtr("0%"),
				Alignment:    &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			})
			if err != nil {
				return err
			}
			if err := f.SetCellStyle(sheet, cell(b.col, dataStart), cell(b.col, dataEnd), pct); err != nil {
				return err
			}
		}
		rng := cell(b.col, dataStart) + ":" + cell(b.col, dataEnd)
		if err := f.SetConditionalFormat(sheet, rng, b.rules); err != nil {
			return err
		}
	}
	return nil
}

// strPtr 返回字符串指针，便于构造 CustomNumFmt 等指针字段。
func strPtr(s string) *string { return &s }
