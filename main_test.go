// main_test.go 包含日历逻辑单元测试与整个工作簿生成的冒烟测试。
package main

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	"hgh/golang_create_excel/demos"
)

// TestBuildCalendar 校验日历排布的正确性。
func TestBuildCalendar(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     time.Month
		wantWeeks int // 期望行数
		wantFirst int // 1 号所在的列（周一=0）
		wantDays  int // 当月天数
	}{
		// 2024-02：闰年 2 月，2 月 1 日是周四（列 3），共 29 天，占 5 行
		{"闰年二月", 2024, time.February, 5, 3, 29},
		// 2023-05：5 月 1 日是周一（列 0），共 31 天，占 5 行
		{"周一起始", 2023, time.May, 5, 0, 31},
		// 2021-02：2 月 1 日是周一（列 0），28 天恰好占满 4 行
		{"恰好四周", 2021, time.February, 4, 0, 28},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weeks := BuildCalendar(tt.year, tt.month)
			if len(weeks) != tt.wantWeeks {
				t.Fatalf("行数 = %d, 期望 %d", len(weeks), tt.wantWeeks)
			}
			if weeks[0][tt.wantFirst] != 1 {
				t.Errorf("1 号应在第 %d 列，实际首行 = %v", tt.wantFirst, weeks[0])
			}
			// 1 号之前必须是空白（0）
			for col := 0; col < tt.wantFirst; col++ {
				if weeks[0][col] != 0 {
					t.Errorf("首行第 %d 列应为空白，实际 = %d", col, weeks[0][col])
				}
			}
			// 全部日期加起来应恰好覆盖 1..daysInMonth
			seen := map[int]bool{}
			for _, w := range weeks {
				for _, d := range w {
					if d != 0 {
						if seen[d] {
							t.Errorf("日期 %d 重复出现", d)
						}
						seen[d] = true
					}
				}
			}
			if len(seen) != tt.wantDays {
				t.Errorf("日期个数 = %d, 期望 %d", len(seen), tt.wantDays)
			}
			if !seen[tt.wantDays] {
				t.Errorf("缺少最后一天 %d", tt.wantDays)
			}
		})
	}
}

// TestWeekdayNames 校验表头顺序为周一到周日。
func TestWeekdayNames(t *testing.T) {
	names := WeekdayNames()
	if names[0] != "周一" || names[6] != "周日" {
		t.Errorf("表头顺序错误: %v", names)
	}
}

// TestCellHelpers 校验单元格地址工具。
func TestCellHelpers(t *testing.T) {
	if got := cell(0, 1); got != "A1" {
		t.Errorf("cell(0, 1) = %q, 期望 %q", got, "A1")
	}
	if got := cell(6, 10); got != "G10" {
		t.Errorf("cell(6, 10) = %q, 期望 %q", got, "G10")
	}
}

// TestRunSmoke 冒烟测试：生成完整工作簿后用 excelize 重新打开校验。
func TestRunSmoke(t *testing.T) {
	out := filepath.Join(t.TempDir(), "showcase.xlsx")
	if err := run(out); err != nil {
		t.Fatalf("run 失败: %v", err)
	}

	f, err := excelize.OpenFile(out)
	if err != nil {
		t.Fatalf("重新打开生成文件失败: %v", err)
	}
	defer f.Close()

	// 1. sheet 齐全且顺序与注册顺序一致
	sheets := f.GetSheetList()
	var want []string
	for _, d := range demos.All() {
		want = append(want, d.Name)
	}
	if len(sheets) != len(want) {
		t.Fatalf("sheet 数量 = %d %v, 期望 %d %v", len(sheets), sheets, len(want), want)
	}
	for i := range want {
		if sheets[i] != want[i] {
			t.Errorf("第 %d 个 sheet = %q, 期望 %q", i, sheets[i], want[i])
		}
	}

	// 2. 月历 sheet 标题应为当前年月
	now := time.Now()
	wantTitle := fmt.Sprintf("%d年%02d月", now.Year(), int(now.Month()))
	got, err := f.GetCellValue("月历", cell(0, 1))
	if err != nil {
		t.Fatalf("读取月历标题失败: %v", err)
	}
	if got != wantTitle {
		t.Errorf("月历标题 = %q, 期望 %q", got, wantTitle)
	}

	// 3. 条件格式 sheet 应存在条件格式规则
	cfs, err := f.GetConditionalFormats("条件格式")
	if err != nil || len(cfs) == 0 {
		t.Errorf("条件格式未写入: cfs=%v err=%v", cfs, err)
	}

	// 4. 图片 sheet 应能取回嵌入的图片
	pics, err := f.GetPictures("图片", "B4")
	if err != nil || len(pics) == 0 {
		t.Errorf("图片未嵌入: pics=%v err=%v", pics, err)
	}

	// 5. 批注 sheet 应存在批注
	comments, err := f.GetComments("批注")
	if err != nil || len(comments) == 0 {
		t.Errorf("批注未写入: comments=%v err=%v", comments, err)
	}
}
