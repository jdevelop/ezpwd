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

func (e *devEzpwd) passwordForm() tview.Primitive {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	onComplete := func(pwd string) {
		crypto, err := ezpwd.NewCrypto([]byte(pwd))
		if err != nil {
			e.showMessage(fmt.Sprintf("can't create crypto: %v", err), screenPwd)
			return
		}
		f, err := os.Open(e.passwordPath)
		switch {
		case err == nil:
			var buf bytes.Buffer
			if err := crypto.Decrypt(f, &buf); err != nil {
				e.showMessage(fmt.Sprintf("can't descrypt storage: %v", err), screenPwd)
				return
			}
			if pwds, err := ezpwd.ReadPasswords(&buf); err != nil {
				e.showMessage(fmt.Sprintf("can't read passwords: %v", err), screenPwd)
				return
			} else {
				e.crypto = crypto
				e.passwordsChan <- pwds
			}
		case errors.Is(err, os.ErrNotExist):
			e.crypto = crypto
			e.passwordsChan <- []ezpwd.Password{
				{
					Service:  "aaa",
					Login:    "bbb",
					Password: "ccc",
					Comment:  "ddd",
				},
			}
		default:
			e.showMessage(fmt.Sprintf("can't open file: %s : %v : %T", e.passwordPath, err, err), screenPwd)
		}
	}
	pwd := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(20).
		SetMaskCharacter('*')
	pwd.
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				onComplete(pwd.GetText())
			case tcell.KeyEsc:
				e.app.Stop()
			}
		})
	form.AddFormItem(pwd).
		AddButton("Unlock", func() {
			onComplete(pwd.GetText())
		}).
		AddButton("Quit", func() {
			e.app.Stop()
		})
	form.SetBorder(true).SetTitle(" Unlock password storage ").SetTitleAlign(tview.AlignCenter)
	form.SetFocus(0)
	return form
}
