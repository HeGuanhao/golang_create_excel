// helpers.go 提供单元格地址相关的通用小工具。
package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// colName 将 0-based 列索引转为 Excel 列名（A, B, C ...）。
func colName(col int) string {
	name, _ := excelize.ColumnNumberToName(col + 1)
	return name
}

// cell 拼接单元格地址，col 为 0-based，row 为 1-based。
// 例如 cell(0, 1) 返回 "A1"。
func cell(col, row int) string {
	return fmt.Sprintf("%s%d", colName(col), row)
}
