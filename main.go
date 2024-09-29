package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"time"
)

func main() {
	// 创建一个新的Excel文件
	f := excelize.NewFile()

	// 获取当前时间
	now := time.Now()
	currentTime := now.Format("2006-01-02 15:04:05") // 格式化时间

	// 在Sheet1的A1单元格写入当前时间
	if err := f.SetCellValue("Sheet1", "A1", currentTime); err != nil {
		fmt.Println(err)
		return
	}

	// 设置文件保存的路径
	fileName := "current_time.xlsx"

	// 保存文件
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Excel文件已生成:", fileName)
}
