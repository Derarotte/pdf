package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
)

type Config struct {
	InputDir    string
	OutputPath  string
	BoxSize     int
	OutputType  string // "svg", "ai", "pdf"
	Spacing     int
}

type PDFExtractor struct {
	config *Config
}

func NewPDFExtractor(config *Config) *PDFExtractor {
	return &PDFExtractor{config: config}
}

// 提取PDF中的矢量图
func (pe *PDFExtractor) ExtractVectorFromPDF(pdfPath string) (image.Image, error) {
	// 使用增强的PDF处理器
	processor := NewEnhancedPDFProcessor(pe.config)
	return processor.ExtractPDFAsImage(pdfPath)
}


// 处理目录中的所有PDF文件
func (pe *PDFExtractor) ProcessDirectory() error {
	// 检查目录是否存在
	if _, err := os.Stat(pe.config.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("输入目录不存在: %s", pe.config.InputDir)
	}

	// 扫描PDF文件，支持大小写不敏感
	var files []string

	// 尝试小写 .pdf
	pdfFiles, err := filepath.Glob(filepath.Join(pe.config.InputDir, "*.pdf"))
	if err == nil {
		files = append(files, pdfFiles...)
	}

	// 尝试大写 .PDF
	PDFFiles, err := filepath.Glob(filepath.Join(pe.config.InputDir, "*.PDF"))
	if err == nil {
		files = append(files, PDFFiles...)
	}

	// 手动扫描目录，以防Glob有问题
	if len(files) == 0 {
		dirEntries, err := os.ReadDir(pe.config.InputDir)
		if err != nil {
			return fmt.Errorf("无法读取目录 %s: %v", pe.config.InputDir, err)
		}

		for _, entry := range dirEntries {
			if !entry.IsDir() {
				name := entry.Name()
				nameLower := strings.ToLower(name)
				if strings.HasSuffix(nameLower, ".pdf") {
					files = append(files, filepath.Join(pe.config.InputDir, name))
				}
			}
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("目录 %s 中没有找到PDF文件\n请检查：\n1. 目录是否包含PDF文件\n2. 文件扩展名是否为.pdf或.PDF", pe.config.InputDir)
	}

	log.Printf("找到 %d 个PDF文件", len(files))

	// 按文件名排序
	sort.Strings(files)

	var extractedImages []image.Image

	for i, file := range files {
		log.Printf("处理文件 (%d/%d): %s", i+1, len(files), filepath.Base(file))

		// 提取矢量图
		img, err := pe.ExtractVectorFromPDF(file)
		if err != nil {
			log.Printf("处理文件 %s 失败: %v", file, err)
			continue
		}

		// 使用增强处理器检测中心点和裁剪
		processor := NewEnhancedPDFProcessor(pe.config)
		center, contentBounds, err := processor.DetectAdvancedBounds(img)
		if err != nil {
			log.Printf("检测中心点失败 %s: %v", file, err)
			continue
		}

		// 智能裁剪图像
		croppedImg := processor.SmartCrop(img, center, contentBounds)
		extractedImages = append(extractedImages, croppedImg)
	}

	if len(extractedImages) == 0 {
		return fmt.Errorf("没有成功提取任何图像")
	}

	// 拼接图像
	return pe.CombineImages(extractedImages)
}

// 拼接图像
func (pe *PDFExtractor) CombineImages(images []image.Image) error {
	if len(images) == 0 {
		return fmt.Errorf("没有图像需要拼接")
	}

	// 计算拼接后的图像尺寸
	totalHeight := len(images) * (pe.config.BoxSize + pe.config.Spacing) - pe.config.Spacing
	combinedImg := image.NewRGBA(image.Rect(0, 0, pe.config.BoxSize, totalHeight))

	// 填充白色背景
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(combinedImg, combinedImg.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// 逐个绘制图像
	currentY := 0
	for _, img := range images {
		dstRect := image.Rect(0, currentY, pe.config.BoxSize, currentY+pe.config.BoxSize)
		draw.Draw(combinedImg, dstRect, img, image.Point{}, draw.Over)
		currentY += pe.config.BoxSize + pe.config.Spacing
	}

	// 保存结果
	return pe.SaveResult(combinedImg)
}

// 保存结果
func (pe *PDFExtractor) SaveResult(img image.Image) error {
	switch pe.config.OutputType {
	case "png":
		return pe.savePNG(img)
	case "svg":
		return pe.saveSVG(img)
	case "ai":
		return pe.saveAI(img)
	case "pdf":
		return pe.savePDF(img)
	default:
		return pe.savePNG(img)
	}
}

func (pe *PDFExtractor) savePNG(img image.Image) error {
	// 生成带时间戳的文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("PDF拼接结果_%s.png", timestamp)
	outputPath := filepath.Join(pe.config.OutputPath, filename)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// 保存为AI格式 (Adobe Illustrator)
func (pe *PDFExtractor) saveAI(img image.Image) error {
	// 生成带时间戳的文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("PDF拼接结果_%s.ai", timestamp)
	outputPath := filepath.Join(pe.config.OutputPath, filename)

	// AI格式实际上是PostScript格式
	// 这里实现一个简化的AI文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 写入AI文件头部
	fmt.Fprintf(file, `%%!PS-Adobe-3.0 EPSF-3.0
%%%%Creator: PDF矢量图提取工具
%%%%BoundingBox: 0 0 %d %d
%%%%DocumentData: Clean7Bit
%%%%LanguageLevel: 2
%%%%Pages: 1
%%%%EndComments
%%%%BeginProlog
%%%%EndProlog
%%%%BeginSetup
%%%%EndSetup
%%%%Page: 1 1
gsave
%d %d scale
/DeviceRGB setcolorspace
`, width, height, width, height)

	// 将图像转换为PostScript数据
	fmt.Fprintf(file, `%d %d 8 [1 0 0 -1 0 %d] {currentfile 3 %d mul string readhexstring pop} false 3 colorimage
`, width, height, height, width)

	// 写入像素数据 (简化版本)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 转换为8位颜色值
			r8 := uint8(r / 257)
			g8 := uint8(g / 257)
			b8 := uint8(b / 257)
			fmt.Fprintf(file, "%02x%02x%02x", r8, g8, b8)

			// 每行添加换行以提高可读性
			if (x-bounds.Min.X+1)%16 == 0 {
				fmt.Fprintf(file, "\n")
			}
		}
	}

	// 写入AI文件尾部
	fmt.Fprintf(file, `
grestore
showpage
%%%%Trailer
%%%%EOF
`)

	return nil
}


func main() {
	// 初始化中文字体支持，解决乱码问题
	initChineseFont()

	myApp := app.New()
	myApp.SetIcon(nil)

	myWindow := myApp.NewWindow("PDF矢量图提取工具 v1.0")
	myWindow.Resize(fyne.NewSize(650, 550))
	myWindow.CenterOnScreen()

	// 配置
	config := &Config{
		BoxSize:    200,
		OutputType: "png",
		Spacing:    10,
	}

	// 界面组件
	inputDirLabel := widget.NewLabel("输入目录:")
	inputDirEntry := widget.NewEntry()
	inputDirEntry.SetPlaceHolder("请选择包含PDF文件的目录路径...")

	// 获取桌面路径
	getDesktopPath := func() string {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(homeDir, "Desktop")
	}

	selectDirBtn := widget.NewButton("选择目录", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if list == nil {
				return
			}
			config.InputDir = list.Path()
			inputDirEntry.SetText(list.Path())
		}, myWindow)
	})

	// 桌面快捷按钮
	desktopBtn := widget.NewButton("桌面", func() {
		desktopPath := getDesktopPath()
		if desktopPath != "" {
			config.InputDir = desktopPath
			inputDirEntry.SetText(desktopPath)
		}
	})

	outputPathLabel := widget.NewLabel("输出文件夹:")
	outputPathEntry := widget.NewEntry()
	outputPathEntry.SetPlaceHolder("请选择输出文件夹，拼接后的图片将保存在此...")

	selectOutputBtn := widget.NewButton("选择输出文件夹", func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if folder == nil {
				return
			}
			config.OutputPath = folder.Path()
			outputPathEntry.SetText(folder.Path())
		}, myWindow)
	})

	// 输出到桌面快捷按钮
	outputDesktopBtn := widget.NewButton("桌面", func() {
		desktopPath := getDesktopPath()
		if desktopPath != "" {
			config.OutputPath = desktopPath
			outputPathEntry.SetText(desktopPath)
		}
	})

	// 参数配置
	boxSizeLabel := widget.NewLabel("裁剪框大小 (像素):")
	boxSizeEntry := widget.NewEntry()
	boxSizeEntry.SetText("200")
	boxSizeEntry.OnChanged = func(text string) {
		if size, err := strconv.Atoi(text); err == nil && size > 0 {
			config.BoxSize = size
		}
	}

	spacingLabel := widget.NewLabel("图片间距 (像素):")
	spacingEntry := widget.NewEntry()
	spacingEntry.SetText("10")
	spacingEntry.OnChanged = func(text string) {
		if spacing, err := strconv.Atoi(text); err == nil && spacing >= 0 {
			config.Spacing = spacing
		}
	}

	outputTypeLabel := widget.NewLabel("输出格式:")
	outputTypeSelect := widget.NewSelect([]string{"PNG", "SVG", "AI", "PDF"}, func(selected string) {
		config.OutputType = strings.ToLower(selected)
	})
	outputTypeSelect.SetSelected("PNG")

	// 状态标签和进度条
	statusLabel := widget.NewLabel("就绪")
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	// 处理按钮
	var processBtn *widget.Button
	processBtn = widget.NewButton("开始处理", func() {
		if config.InputDir == "" {
			dialog.ShowError(fmt.Errorf("请选择输入目录"), myWindow)
			return
		}
		if config.OutputPath == "" {
			dialog.ShowError(fmt.Errorf("请选择输出文件夹"), myWindow)
			return
		}

		// 异步处理，避免界面卡顿
		processBtn.SetText("处理中...")
		processBtn.Disable()
		statusLabel.SetText("正在处理PDF文件...")
		progressBar.Show()
		progressBar.Start()

		go func() {
			extractor := NewPDFExtractor(config)

			// 首先检查目录和文件
			statusLabel.SetText("正在扫描PDF文件...")

			err := extractor.ProcessDirectory()

			// 在主线程更新UI
			processBtn.SetText("开始处理")
			processBtn.Enable()
			progressBar.Stop()
			progressBar.Hide()

			if err != nil {
				statusLabel.SetText("处理失败")
				dialog.ShowError(err, myWindow)
			} else {
				statusLabel.SetText("处理完成")
				dialog.ShowInformation("成功", "拼接完成！文件已保存到输出文件夹。", myWindow)
			}
		}()
	})

	// 布局
	content := container.NewVBox(
		inputDirLabel,
		container.NewBorder(nil, nil, desktopBtn, selectDirBtn, inputDirEntry),
		widget.NewSeparator(),

		outputPathLabel,
		container.NewBorder(nil, nil, outputDesktopBtn, selectOutputBtn, outputPathEntry),
		widget.NewSeparator(),

		container.NewGridWithColumns(2,
			boxSizeLabel, boxSizeEntry,
			spacingLabel, spacingEntry,
			outputTypeLabel, outputTypeSelect,
		),
		widget.NewSeparator(),

		statusLabel,
		progressBar,
		processBtn,
	)

	myWindow.SetContent(container.NewPadded(content))
	myWindow.ShowAndRun()
}