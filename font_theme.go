package main

import (
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// 中文字体主题
type chineseTheme struct{}

func newChineseTheme() fyne.Theme {
	return &chineseTheme{}
}

func (ct *chineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (ct *chineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 对于Windows，Fyne会自动使用系统字体
	// 但我们可以明确指定使用默认字体以确保中文支持
	return theme.DefaultTheme().Font(style)
}

func (ct *chineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (ct *chineseTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}