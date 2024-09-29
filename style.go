package main

import "github.com/xuri/excelize/v2"

// 定义一个结构体来存储样式ID和可能的额外信息
type ExcelStyle struct {
	ID        int
	FillColor []uint8 // 示例：仅包含填充色
	// 可以添加更多字段来表示不同的样式属性
}

// 创建一个新的灰色背景样式
func CreateGrayBackgroundStyle(f *excelize.File) (*ExcelStyle, error) {
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"10526880"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	return &ExcelStyle{ID: style, FillColor: []uint8{217, 217, 217}}, nil
}

// 创建一个新的加粗样式
func CreateBoldStyle(f *excelize.File) (*ExcelStyle, error) {
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return &ExcelStyle{ID: style}, nil
}
