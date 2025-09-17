package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
)

// SVG输出模块
func (pe *PDFExtractor) saveSVG(img image.Image) error {
	outputPath := pe.config.OutputPath
	if !strings.HasSuffix(outputPath, ".svg") {
		outputPath += ".svg"
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建SVG文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入SVG头部
	fmt.Fprintf(file, `<?xml version="1.0" encoding="UTF-8"?>
<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
`, width, height, width, height)

	// 将图像转换为SVG矩形
	// 这是一个简化的实现，实际项目中可能需要更复杂的矢量化算法
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 { // 降采样以减小文件大小
		for x := bounds.Min.X; x < bounds.Max.X; x += 2 {
			r, g, b, a := img.At(x, y).RGBA()
			if a > 0 && (r < 65000 || g < 65000 || b < 65000) { // 不是白色
				// 转换为8位颜色值
				r8 := uint8(r / 257)
				g8 := uint8(g / 257)
				b8 := uint8(b / 257)

				fmt.Fprintf(file, `  <rect x="%d" y="%d" width="2" height="2" fill="#%02x%02x%02x"/>
`, x-bounds.Min.X, y-bounds.Min.Y, r8, g8, b8)
			}
		}
	}

	// 写入SVG尾部
	fmt.Fprintf(file, "</svg>\n")

	return nil
}

// 改进的SVG输出，使用路径而不是矩形
func (pe *PDFExtractor) saveSVGOptimized(img image.Image) error {
	outputPath := pe.config.OutputPath
	if !strings.HasSuffix(outputPath, ".svg") {
		outputPath += ".svg"
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入SVG头部
	fmt.Fprintf(file, `<?xml version="1.0" encoding="UTF-8"?>
<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
`, width, height, width, height)

	// 按颜色分组像素，生成优化的路径
	colorPaths := make(map[color.RGBA][]string)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a > 0 && (r < 65000 || g < 65000 || b < 65000) {
				c := color.RGBA{
					R: uint8(r / 257),
					G: uint8(g / 257),
					B: uint8(b / 257),
					A: 255,
				}

				path := fmt.Sprintf("M%d,%dh1v1h-1z", x-bounds.Min.X, y-bounds.Min.Y)
				colorPaths[c] = append(colorPaths[c], path)
			}
		}
	}

	// 为每种颜色生成路径
	for c, paths := range colorPaths {
		if len(paths) > 0 {
			fmt.Fprintf(file, `  <path d="%s" fill="#%02x%02x%02x"/>
`, strings.Join(paths, ""), c.R, c.G, c.B)
		}
	}

	fmt.Fprintf(file, "</svg>\n")
	return nil
}