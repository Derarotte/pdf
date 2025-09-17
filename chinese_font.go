package main

import (
	"os"
	"strings"

	"github.com/flopp/go-findfont"
)

// 初始化中文字体支持，解决中文乱码问题
func initChineseFont() {
	// 获取系统所有字体路径
	fontPaths := findfont.List()

	// 常见的中文字体文件名
	chineseFonts := []string{
		"msyh.ttf",     // 微软雅黑
		"msyhbd.ttf",   // 微软雅黑粗体
		"simhei.ttf",   // 黑体
		"simsun.ttc",   // 宋体
		"simkai.ttf",   // 楷体
		"simfang.ttf",  // 仿宋
		"simli.ttf",    // 隶书
		"simyou.ttf",   // 幼圆
		"STXIHEI.TTF",  // 华文细黑
		"STKAITI.TTF",  // 华文楷体
		"STFANGSO.TTF", // 华文仿宋
		"STZHONGS.TTF", // 华文中宋
		"NotoSansCJK-Regular.ttc", // Noto CJK字体
		"DengXian.ttf", // 等线体
		"YaHei.ttf",    // 雅黑
	}

	// 查找并设置中文字体
	for _, path := range fontPaths {
		pathLower := strings.ToLower(path)
		for _, chineseFont := range chineseFonts {
			if strings.Contains(pathLower, strings.ToLower(chineseFont)) {
				// 设置FYNE_FONT环境变量
				os.Setenv("FYNE_FONT", path)
				return // 找到第一个可用的中文字体就返回
			}
		}
	}

	// 如果没有找到特定的中文字体，尝试寻找任何包含"chinese"、"cjk"、"han"的字体
	fallbackKeywords := []string{"chinese", "cjk", "han", "zh", "cn"}
	for _, path := range fontPaths {
		pathLower := strings.ToLower(path)
		for _, keyword := range fallbackKeywords {
			if strings.Contains(pathLower, keyword) {
				os.Setenv("FYNE_FONT", path)
				return
			}
		}
	}
}