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

	// 调用styles.go中定义的函数来创建样式
	grayStyle, err := CreateGrayBackgroundStyle(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	boldStyle, err := CreateBoldStyle(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 应用样式
	f.SetColWidth("events", "A", "A", 20)
	f.SetColWidth("events", "B", "B", 200)
	f.SetCellValue("Sheet1", "A1", "标题")
	f.SetCellStyle("Sheet1", "A1", "A1", grayStyle.ID)
	f.SetCellValue("Sheet1", "B1", "时间")
	f.SetCellStyle("Sheet1", "B1", "B1", grayStyle.ID)

	f.SetCellValue("Sheet1", "B2", currentTime)
	f.SetCellStyle("Sheet1", "B2", "B2", boldStyle.ID)

	// 设置文件保存的路径
	fileName := "current_time.xlsx"

	// 保存文件
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Excel文件已生成:", fileName)
}
