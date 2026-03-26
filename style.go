package main

import "github.com/xuri/excelize/v2"

type Styles struct {
	Title   int // 年月标题
	Header  int // 星期表头
	Weekend int // 周末
	Weekday int // 工作日
	Today   int // 今天
	Empty   int // 空白格（非本月）
	InfoKey int // 信息栏 key
	InfoVal int // 信息栏 value
}

func NewStyles(f *excelize.File) (*Styles, error) {
	title, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2E75B6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	header, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "top", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	weekend, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "C00000"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"FFE7E7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "DDDDDD", Style: 1},
			{Type: "right", Color: "DDDDDD", Style: 1},
			{Type: "top", Color: "DDDDDD", Style: 1},
			{Type: "bottom", Color: "DDDDDD", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	weekday, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Color: "333333"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F2F2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "DDDDDD", Style: 1},
			{Type: "right", Color: "DDDDDD", Style: 1},
			{Type: "top", Color: "DDDDDD", Style: 1},
			{Type: "bottom", Color: "DDDDDD", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	today, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"ED7D31"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "C55A11", Style: 2},
			{Type: "right", Color: "C55A11", Style: 2},
			{Type: "top", Color: "C55A11", Style: 2},
			{Type: "bottom", Color: "C55A11", Style: 2},
		},
	})
	if err != nil {
		return nil, err
	}

	empty, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "EEEEEE", Style: 1},
			{Type: "right", Color: "EEEEEE", Style: 1},
			{Type: "top", Color: "EEEEEE", Style: 1},
			{Type: "bottom", Color: "EEEEEE", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	infoKey, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	infoVal, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "333333"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"EBF3FB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	return &Styles{
		Title:   title,
		Header:  header,
		Weekend: weekend,
		Weekday: weekday,
		Today:   today,
		Empty:   empty,
		InfoKey: infoKey,
		InfoVal: infoVal,
	}, nil
}
