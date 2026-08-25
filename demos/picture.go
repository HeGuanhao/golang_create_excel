// picture.go 渲染"图片" sheet。
//
// 图片不依赖任何外部文件：用标准库 image/png 在内存中现场绘制，
// 再通过 AddPictureFromBytes 嵌入工作簿。
package demos

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/xuri/excelize/v2"
)

// genBannerPNG 生成一张 360x200 的渐变背景 + 色块 PNG。
func genBannerPNG() ([]byte, error) {
	const w, h = 360, 200
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// 纵向渐变背景：主题蓝 -> 浅蓝
	for y := 0; y < h; y++ {
		c := color.RGBA{
			R: uint8(46 + y*100/h),
			G: uint8(117 + y*60/h),
			B: uint8(182 + y*40/h),
			A: 255,
		}
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	// 白色内边框
	for i := 0; i < 4; i++ {
		for x := i; x < w-i; x++ {
			img.Set(x, i, color.White)
			img.Set(x, h-1-i, color.White)
		}
		for y := i; y < h-i; y++ {
			img.Set(i, y, color.White)
			img.Set(w-1-i, y, color.White)
		}
	}
	// 底部三个色块（橙 / 红 / 深灰），模拟简单图形
	blocks := []color.RGBA{
		{237, 125, 49, 255}, // 强调橙
		{192, 0, 0, 255},    // 周末红
		{51, 51, 51, 255},   // 正文灰
	}
	for i, bc := range blocks {
		x0 := 40 + i*100
		for y := 120; y < 170; y++ {
			for x := x0; x < x0+60; x++ {
				img.Set(x, y, bc)
			}
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// genCheckerPNG 生成一张 96x96 的棋盘格 PNG。
func genCheckerPNG() ([]byte, error) {
	const size, cellPx = 96, 16
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	c1 := color.RGBA{68, 114, 196, 255}  // 表头蓝
	c2 := color.RGBA{242, 242, 242, 255} // 浅灰
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if (x/cellPx+y/cellPx)%2 == 0 {
				img.Set(x, y, c1)
			} else {
				img.Set(x, y, c2)
			}
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func registerPicture(f *excelize.File) error {
	const sheet = "图片"
	if err := newSheet(f, sheet, "下方两张图均由 image/png 在内存中现场绘制，未使用任何外部图片文件", 8); err != nil {
		return err
	}

	banner, err := genBannerPNG()
	if err != nil {
		return err
	}
	checker, err := genCheckerPNG()
	if err != nil {
		return err
	}

	// 主图：渐变横幅
	if err := f.AddPictureFromBytes(sheet, "B4", &excelize.Picture{
		Extension: ".png",
		File:      banner,
		Format:    &excelize.GraphicOptions{OffsetX: 10, OffsetY: 10},
	}); err != nil {
		return err
	}
	// 小图：棋盘格，放在横幅右侧
	if err := f.AddPictureFromBytes(sheet, "I4", &excelize.Picture{
		Extension: ".png",
		File:      checker,
		Format: &excelize.GraphicOptions{
			OffsetX: 10, OffsetY: 10,
			ScaleX: 1.5, ScaleY: 1.5,
		},
	}); err != nil {
		return err
	}
	return nil
}
