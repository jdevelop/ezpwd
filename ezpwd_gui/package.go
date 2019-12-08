package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func (e *devEzpwd) showMessage(msg string, previousScreen string) {
	e.app.QueueUpdateDraw(func() {
		text := tview.NewTextView().
			SetText(msg).
			SetWrap(true)
		text.
			SetTitle(" Error ").
			SetTitleColor(tcell.ColorRed).
			SetBorder(true).
			SetBorderColor(tcell.ColorRed).
			SetBorderPadding(1, 1, 1, 1)
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 2, false).
			AddItem(
				tview.NewFlex().
					AddItem(nil, 0, 2, false).
					AddItem(text, len(msg)+6, 0, true).
					AddItem(nil, 0, 2, false),
				6, 0, true,
			).
			AddItem(nil, 0, 2, false)
		e.pages.AddPage(screenError, flex, true, true)
		text.SetDoneFunc(func(k tcell.Key) {
			switch k {
			case tcell.KeyEnter, tcell.KeyEsc:
				e.pages.RemovePage(screenError)
				e.pages.SwitchToPage(previousScreen)
			}
		})
		e.app.SetFocus(text)
	})
}

func (e *devEzpwd) confirm(msg, from string, ok func()) {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	form.SetTitle(msg)
	form.SetBorder(true)
	form.AddButton("Ok", func() {
		ok()
		e.pages.RemovePage(screenConfirm)
		e.pages.ShowPage(from)
	}).AddButton("Cancel", func() {
		e.pages.RemovePage(screenConfirm)
		e.pages.ShowPage(from)
	})
	e.pages.AddPage(screenConfirm, modal(form, len(msg)+4, 5), true, true)
	e.app.SetFocus(form)
}
