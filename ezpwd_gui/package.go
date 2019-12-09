package main

import (
	"fmt"
	"io"
	"os"
	"time"

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

func (e *devEzpwd) showMessage(title, msg, previousScreen string, styleUpd ...func(*tview.TextView)) {
	e.app.QueueUpdateDraw(func() {
		text := tview.NewTextView().
			SetText(msg).
			SetWrap(true).
			SetTextAlign(tview.AlignCenter)
		text.
			SetTitle(title).
			SetTitleColor(tcell.ColorRed).
			SetBorder(true).
			SetBorderColor(tcell.ColorRed).
			SetBorderPadding(1, 1, 1, 1)
		for _, f := range styleUpd {
			f(text)
		}
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

const layout = "020106"

func backup(file string) error {
	current := time.Now()
	newName := fmt.Sprintf("%s.%s-%d", file, current.Format(layout), current.Unix())
	src, err := os.Open(file)
	if err != nil {
		return err
	}
	f, err := os.Create(newName)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, src)
	return err
}
