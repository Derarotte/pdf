package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

// 简化的PDF输出模块
func (pe *PDFExtractor) savePDF(img image.Image) error {
	// 生成带时间戳的文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("PDF拼接结果_%s.pdf", timestamp)
	outputPath := filepath.Join(pe.config.OutputPath, filename)

	// 这里使用一个简化的PDF生成
	// 实际项目中建议使用专业的PDF生成库如 gofpdf
	return pe.generateSimplePDF(img, outputPath)
}

// 生成简单的PDF文件
func (pe *PDFExtractor) generateSimplePDF(img image.Image, outputPath string) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 将图像编码为PNG字节
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return err
	}
	imageData := buf.Bytes()

	// 创建PDF文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 生成简化的PDF内容
	// 这是一个最基本的PDF结构，包含一个图像对象
	pdfContent := fmt.Sprintf(`%%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj

2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj

3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 %d %d]
/Resources <<
  /XObject << /Im1 4 0 R >>
>>
/Contents 5 0 R
>>
endobj

4 0 obj
<<
/Type /XObject
/Subtype /Image
/Width %d
/Height %d
/ColorSpace /DeviceRGB
/BitsPerComponent 8
/Filter /FlateDecode
/Length %d
>>
stream
`, width, height, width, height, len(imageData))

	// 写入PDF内容
	_, err = file.WriteString(pdfContent)
	if err != nil {
		return err
	}

	// 写入图像数据
	_, err = file.Write(imageData)
	if err != nil {
		return err
	}

	// 完成PDF结构
	endContent := fmt.Sprintf(`
endstream
endobj

5 0 obj
<<
/Length 44
>>
stream
q
%d 0 0 %d 0 0 cm
/Im1 Do
Q
endstream
endobj

xref
0 6
0000000000 65535 f
0000000010 00000 n
0000000079 00000 n
0000000136 00000 n
0000000301 00000 n
0000000500 00000 n
trailer
<<
/Size 6
/Root 1 0 R
>>
startxref
600
%%%%EOF
`, width, height)

	_, err = file.WriteString(endContent)
	return err
}