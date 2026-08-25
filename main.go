// main.go 程序入口：编排所有 demo，输出 showcase.xlsx。
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"

	"hgh/golang_create_excel/demos"
)

// outputFile 生成的展示工作簿文件名。
const outputFile = "showcase.xlsx"

func main() {
	if err := run(outputFile); err != nil {
		fmt.Fprintln(os.Stderr, "生成失败:", err)
		os.Exit(1)
	}
}

// run 生成特效展示工作簿并保存到 out。
func run(out string) error {
	f := excelize.NewFile()
	defer f.Close()

	// 逐个渲染 demo sheet，任一失败则中止并给出明确错误
	for _, d := range demos.All() {
		if err := d.Register(f); err != nil {
			return fmt.Errorf("渲染工作表 %q（%s）失败: %w", d.Name, d.Desc, err)
		}
	}
	// 删除 excelize 默认创建的 Sheet1，并激活第一个 sheet
	if err := f.DeleteSheet("Sheet1"); err != nil {
		return fmt.Errorf("删除默认工作表失败: %w", err)
	}
	if idx, err := f.GetSheetIndex(demos.All()[0].Name); err == nil {
		f.SetActiveSheet(idx)
	}

	if err := f.SaveAs(out); err != nil {
		return fmt.Errorf("保存文件 %s 失败: %w", out, err)
	}

	abs, err := filepath.Abs(out)
	if err != nil {
		abs = out
	}
	fmt.Println("生成成功:", abs)
	fmt.Println("包含的工作表:")
	for i, name := range f.GetSheetList() {
		fmt.Printf("  %d. %s\n", i+1, name)
	}
	return nil
}
