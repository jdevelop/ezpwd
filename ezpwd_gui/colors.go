package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	BackgroundColor = tcell.NewRGBColor(0x87, 0x87, 0x5f)
	FormBolderColor = tcell.NewRGBColor(0xff, 0xf6, 0xe9)
)

func init() {
	tview.Styles.PrimitiveBackgroundColor = BackgroundColor
	tview.Styles.BorderColor = FormBolderColor
	tview.Styles.SecondaryTextColor = tcell.ColorBlack
	tview.Styles.TitleColor = tcell.ColorBlack
	tview.Styles.GraphicsColor = tcell.ColorBlack
}
