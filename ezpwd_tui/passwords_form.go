package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

func (e *devEzpwd) initForm() *tview.Form {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	initFormColors(form)
	pwd := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(20).
		SetMaskCharacter('*')
	confirm := tview.NewInputField().
		SetLabel("Confirm").
		SetFieldWidth(20).
		SetMaskCharacter('*')
	onComplete := func() {
		if pwd.GetText() == "" {
			e.showMessage("Error", "Password can't be empty", screenPwd, errorsMessageStyle)
			return
		}
		if pwd.GetText() != confirm.GetText() {
			e.showMessage("Error", "Password doesn't match confirmation", screenPwd, errorsMessageStyle)
			return
		}
		crypto, err := ezpwd.NewCrypto([]byte(pwd.GetText()))
		if err != nil {
			e.showMessage("Error", fmt.Sprintf("can't create crypto: %v", err), screenPwd, errorsMessageStyle)
			return
		}
		e.crypto = crypto
		e.passwordsChan <- []ezpwd.Password{}
	}
	form.SetCancelFunc(e.app.Stop)
	form.
		AddFormItem(pwd).
		AddFormItem(confirm).
		AddButton("Create", func() {
			onComplete()
		}).
		AddButton("Quit", func() {
			e.app.Stop()
		})
	form.SetBorder(true).SetTitle(" Create password storage ").SetTitleAlign(tview.AlignCenter)
	form.SetFocus(0)
	return form
}

func initFormColors(form *tview.Form) {
	form.SetBackgroundColor(DefaultColorSchema.LoginFormColors.Background)
	form.SetTitleColor(DefaultColorSchema.LoginFormColors.Title)
	form.SetBorderColor(DefaultColorSchema.LoginFormColors.Border)
	form.SetLabelColor(DefaultColorSchema.LoginFormColors.Label)
	form.SetButtonBackgroundColor(DefaultColorSchema.LoginFormColors.ButtonBackground)
	form.SetButtonTextColor(DefaultColorSchema.LoginFormColors.ButtonText)
	form.SetFieldBackgroundColor(DefaultColorSchema.LoginFormColors.FieldBackground)
	form.SetFieldTextColor(DefaultColorSchema.LoginFormColors.FieldText)
}

func (e *devEzpwd) passwordForm() *tview.Form {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	initFormColors(form)
	onComplete := func(pwd string) {
		crypto, err := ezpwd.NewCrypto([]byte(pwd))
		if err != nil {
			e.showMessage("Error", fmt.Sprintf("can't create crypto: %v", err), screenPwd, errorsMessageStyle)
			return
		}
		f, err := os.Open(e.passwordPath)
		switch {
		case err == nil:
			var buf bytes.Buffer
			if err := crypto.Decrypt(f, &buf); err != nil {
				e.showMessage("Error", fmt.Sprintf("can't descrypt storage: %v", err), screenPwd, errorsMessageStyle)
				return
			}
			if pwds, err := ezpwd.ReadPasswords(&buf); err != nil {
				e.showMessage("Error", fmt.Sprintf("can't read passwords: %v", err), screenPwd, errorsMessageStyle)
				return
			} else {
				e.crypto = crypto
				e.passwordsChan <- pwds
			}
		case errors.Is(err, os.ErrNotExist):
			e.crypto = crypto
			e.passwordsChan <- []ezpwd.Password{}
		default:
			e.showMessage("Error", fmt.Sprintf("can't open file: %s : %v : %T", e.passwordPath, err, err), screenPwd, errorsMessageStyle)
		}
	}
	pwd := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(20).
		SetMaskCharacter('*')

	escHandler := func() {
		if pwd.GetText() == "" {
			e.app.Stop()
		} else {
			pwd.SetText("")
			e.app.SetFocus(pwd)
		}
	}
	pwd.
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				onComplete(pwd.GetText())
			case tcell.KeyEsc:
				escHandler()
			}
		})
	form.SetCancelFunc(escHandler)
	btnOk := tview.NewButton("Unlock").SetSelectedFunc(func() { onComplete(pwd.GetText()) })
	btnOk.SetBackgroundColor(tcell.ColorRed)
	form.AddFormItem(pwd).
		AddButton("Unlock", func() { onComplete(pwd.GetText()) }).
		AddButton("Quit", e.app.Stop)
	form.SetBorder(true).SetTitle(" Unlock password storage ").SetTitleAlign(tview.AlignCenter)
	form.SetFocus(0)
	return form
}
