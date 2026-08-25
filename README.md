# golang_create_excel

使用 Go 与 [excelize v2.8.1](https://github.com/xuri/excelize) 生成 Excel 的**特效示例集**。运行一次即可生成一个 `showcase.xlsx`，其中每个工作表演示一种渲染效果。

## 运行方式

```bash
go run .        # 在当前目录生成 showcase.xlsx
go test ./...   # 运行日历逻辑单测与工作簿生成冒烟测试
```

## 效果清单

| 工作表 | 演示效果 | 源文件 |
| --- | --- | --- |
| 月历 | 当月日历：年月标题栏、周一到周日表头、周末红色高亮、今天橙色高亮、底部信息栏 | `demos/calendar.go`（逻辑） |
| 样式一览 | 字体（加粗/斜体/下划线/删除线/颜色/字号/字体族）、填充、边框（细/中/虚线）、对齐（左中右/自动换行/旋转）、数字格式（千分位/百分比/货币/日期/科学计数） | `demos/styles.go` |
| 条件格式 | 数据条、三阶色阶、交通灯图标集、">90 高亮"单元格规则 | `demos/conditional.go` |
| 图表 | 柱状图、折线图（双序列带数据标记）、饼图（百分比数据标签） | `demos/chart.go` |
| 表格 | `AddTable` 表格样式（TableStyleMedium9）、`SetPanes` 冻结首行、普通区域 `AutoFilter` 自动筛选 | `demos/table.go` |
| 富文本 | `SetCellRichText`：同一单元格内多字体/颜色混排、上下标（H₂O、E=mc²） | `demos/richtext.go` |
| 数据验证 | 下拉列表（部门）、整数范围校验（年龄 18-65），含输入提示与出错警告 | `demos/validation.go` |
| 图片 | 用标准库 `image/png` 在内存中现场绘制 PNG（渐变横幅、棋盘格），经 `AddPictureFromBytes` 嵌入，不依赖任何外部图片文件 | `demos/picture.go` |
| 批注 | 普通文本批注与富文本批注 | `demos/comment.go` |

## 项目结构

```
├── main.go            # 入口：编排所有 demo，输出 showcase.xlsx，任一失败返回非零退出码
├── helpers.go         # 单元格地址小工具（colName / cell）
├── calendar.go        # 日历纯逻辑 API（转发到 demos 包，保持可测试）
├── demos/
│   ├── demo.go        # Demo 注册接口 + 主题色板与样式小工具
│   ├── calendar.go    # 日历逻辑实现与"月历" sheet 渲染
│   ├── styles.go      # 样式一览
│   ├── conditional.go # 条件格式
│   ├── chart.go       # 图表
│   ├── table.go       # 表格 / 冻结窗格 / 自动筛选
│   ├── richtext.go    # 富文本
│   ├── validation.go  # 数据验证
│   ├── picture.go     # 内存生成 PNG 并嵌入
│   └── comment.go     # 批注
├── main_test.go       # 日历逻辑单测 + 生成后重新打开校验的冒烟测试
└── README.md
```

> 说明：Go 不允许子包导入 main 包，因此日历逻辑的实现放在 `demos/calendar.go`，
> 根目录 `calendar.go` 以类型别名 + 转发函数的形式保留原有 API，供根包测试使用。

## 设计要点

- **统一视觉**：`demos/demo.go` 中的主题色板（标题蓝、表头蓝、强调橙、周末红）与各样式小工具（`fill`、`borders`、`addStyle`、`newSheet`）被所有 demo 复用，避免重复的样式代码；每个 sheet 顶部都有统一的标题栏与说明栏。
- **错误处理**：任一 demo 渲染失败都会打印明确错误（指明是哪个 sheet）并以非零退出码退出。
- **可测试**：日历排布为纯函数，直接单元测试；整个工作簿的生成有冒烟测试，生成后用 `excelize.OpenFile` 重新打开校验 sheet 列表、月历标题、条件格式、图片与批注。

## 效果截图

生成 `showcase.xlsx` 后用 Excel / WPS 打开即可查看各 sheet 效果（截图占位）：

- 月历：`docs/screenshots/calendar.png`
- 样式一览：`docs/screenshots/styles.png`
- 条件格式：`docs/screenshots/conditional.png`
- 图表：`docs/screenshots/chart.png`
