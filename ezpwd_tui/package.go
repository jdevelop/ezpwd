package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func modal(p tview.Primitive, width, height int, style ...func(p *tview.Box)) tview.Primitive {
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
	for _, f := range style {
		f(flex.Box)
	}
	return flex
}

func (e *devEzpwd) showMessage(title, msg, previousScreen string, styleUpd ...func(*tview.TextView)) {
	e.app.QueueUpdateDraw(func() {
		text := tview.NewTextView().
			SetText(msg).
			SetWrap(true).
			SetTextAlign(tview.AlignCenter)
		text.
			SetTitle(title).
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1)
		for _, f := range styleUpd {
			f(text)
		}
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(
				tview.NewFlex().
					AddItem(nil, 0, 1, false).
					AddItem(text, len(msg)+6, 0, true).
					AddItem(nil, 0, 1, false),
				5, 0, true,
			).
			AddItem(nil, 0, 1, false)
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
	form.SetBackgroundColor(DefaultColorSchema.ConfirmFormColors.Background)
	form.SetLabelColor(DefaultColorSchema.ConfirmFormColors.Label)
	form.SetButtonBackgroundColor(DefaultColorSchema.ConfirmFormColors.ButtonBackground)
	form.SetButtonTextColor(DefaultColorSchema.ConfirmFormColors.ButtonText)
	form.SetFieldBackgroundColor(DefaultColorSchema.ConfirmFormColors.FieldBackground)
	form.SetFieldTextColor(DefaultColorSchema.ConfirmFormColors.FieldText)
	cancelFunc := func() {
		e.pages.RemovePage(screenConfirm)
		e.pages.ShowPage(from)
	}
	form.SetCancelFunc(cancelFunc)
	form.AddButton("Ok", func() {
		ok()
		e.pages.RemovePage(screenConfirm)
		e.pages.ShowPage(from)
	}).AddButton("Cancel", cancelFunc)
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

func successMessageStyle(text *tview.TextView) {
	text.SetBorderColor(DefaultColorSchema.MessagesColors.SuccessBorder)
	text.SetTitleColor(DefaultColorSchema.MessagesColors.SuccessTitle)
	text.SetBackgroundColor(DefaultColorSchema.MessagesColors.SuccessBackground)
	text.SetTextColor(DefaultColorSchema.MessagesColors.SuccessText)
}

func errorsMessageStyle(text *tview.TextView) {
	text.SetBorderColor(DefaultColorSchema.MessagesColors.FailureBorder)
	text.SetTitleColor(DefaultColorSchema.MessagesColors.FailureTitle)
	text.SetBackgroundColor(DefaultColorSchema.MessagesColors.FailureBackground)
	text.SetTextColor(DefaultColorSchema.MessagesColors.FailureText)
}
