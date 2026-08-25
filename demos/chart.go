// chart.go 渲染"图表" sheet：柱状图、折线图、饼图。
//
// v2.8.1 中 AddChart 的签名为 AddChart(sheet, cell string, chart *Chart, combo ...*Chart)，
// 图表数据直接引用本 sheet 左侧的数据区（中文 sheet 名引用需加单引号）。
package demos

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func registerChart(f *excelize.File) error {
	const sheet = "图表"
	if err := newSheet(f, sheet, "A 列为图表数据源；右侧依次为柱状图、折线图、饼图", 9); err != nil {
		return err
	}

	// ── 数据源 ──
	// 月度数据：A4:C10
	months := [][]interface{}{
		{"月份", "销量(台)", "订单量"},
		{"1月", 120, 98},
		{"2月", 96, 110},
		{"3月", 158, 132},
		{"4月", 132, 121},
		{"5月", 171, 160},
		{"6月", 145, 149},
	}
	if err := setRows(f, sheet, 1, 4, months); err != nil {
		return err
	}
	// 渠道占比：A13:B17（饼图数据）
	channels := [][]interface{}{
		{"渠道", "占比"},
		{"线上", 45},
		{"门店", 30},
		{"代理", 15},
		{"其他", 10},
	}
	if err := setRows(f, sheet, 1, 13, channels); err != nil {
		return err
	}

	headerStyle, err := addStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"},
		Fill:      fill(colorHeader),
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}
	for _, rng := range [][2]string{{"A4", "C4"}, {"A13", "B13"}} {
		if err := f.SetCellStyle(sheet, rng[0], rng[1], headerStyle); err != nil {
			return err
		}
	}
	if err := f.SetColWidth(sheet, "A", "C", 12); err != nil {
		return err
	}

	// 引用本 sheet 数据（中文 sheet 名需单引号包裹）
	ref := func(rng string) string { return fmt.Sprintf("'%s'!%s", sheet, rng) }
	title := func(s string) []excelize.RichTextRun { return []excelize.RichTextRun{{Text: s}} }
	dim := excelize.ChartDimension{Width: 520, Height: 300}

	// 柱状图：各月销量
	if err := f.AddChart(sheet, "E4", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{{
			Name:       ref("$B$4"),
			Categories: ref("$A$5:$A$10"),
			Values:     ref("$B$5:$B$10"),
		}},
		Title:     title("2024 上半年销量柱状图"),
		Dimension: dim,
		Legend:    excelize.ChartLegend{Position: "bottom"},
	}); err != nil {
		return err
	}

	// 折线图：销量 vs 订单量
	if err := f.AddChart(sheet, "E20", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       ref("$B$4"),
				Categories: ref("$A$5:$A$10"),
				Values:     ref("$B$5:$B$10"),
				Marker:     excelize.ChartMarker{Symbol: "circle", Size: 6},
			},
			{
				Name:       ref("$C$4"),
				Categories: ref("$A$5:$A$10"),
				Values:     ref("$C$5:$C$10"),
				Marker:     excelize.ChartMarker{Symbol: "diamond", Size: 6},
			},
		},
		Title:     title("销量与订单量走势"),
		Dimension: dim,
		Legend:    excelize.ChartLegend{Position: "bottom"},
	}); err != nil {
		return err
	}

	// 饼图：渠道占比
	if err := f.AddChart(sheet, "M4", &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{{
			Name:       ref("$B$13"),
			Categories: ref("$A$14:$A$17"),
			Values:     ref("$B$14:$B$17"),
		}},
		Title:     title("销售渠道占比"),
		Dimension: dim,
		Legend:    excelize.ChartLegend{Position: "right"},
		PlotArea: excelize.ChartPlotArea{
			ShowPercent: true, // 数据标签显示百分比
		},
	}); err != nil {
		return err
	}
	return nil
}
