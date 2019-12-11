package main

import (
	"fmt"
	"math/rand"

	"github.com/dchest/uniuri"
	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

func (e *devEzpwd) passwordMgmtForm(id int, pwds []ezpwd.Password) *tview.Form {
	var genPwdFunc func()
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	form.SetBackgroundColor(passwordMgmgFormColors.Background)
	form.SetLabelColor(passwordMgmgFormColors.Label)
	form.SetButtonBackgroundColor(passwordMgmgFormColors.ButtonBackground)
	form.SetButtonTextColor(passwordMgmgFormColors.ButtonText)
	form.SetFieldBackgroundColor(passwordMgmgFormColors.FieldBackground)
	form.SetFieldTextColor(passwordMgmgFormColors.FieldText)
	svc := tview.NewInputField().SetLabel("Service:").SetFieldWidth(40)
	login := tview.NewInputField().SetLabel("Login:").SetFieldWidth(40)
	comment := tview.NewInputField().SetLabel("Comment:").SetFieldWidth(40)
	setupPwgGen := func(i *tview.InputField) *tview.InputField {
		i.SetPlaceholder("Press 'Alt-G' to generate").
			SetInputCapture(func(k *tcell.EventKey) *tcell.EventKey {
				if k.Modifiers()&tcell.ModAlt > 0 && k.Rune() == 'g' {
					genPwdFunc()
					form.SetFocus(4)
					e.app.SetFocus(form)
					return nil
				} else {
					return k
				}
			})
		return i
	}
	pwd := setupPwgGen(tview.NewInputField().SetLabel("Password:").SetFieldWidth(40).
		SetMaskCharacter('*'))
	confirm := setupPwgGen(tview.NewInputField().SetLabel("Confirm:").SetFieldWidth(40).
		SetMaskCharacter('*'))
	genPwdFunc = func() {
		password := uniuri.NewLen(rand.Intn(8) + 8)
		pwd.SetText(password)
		confirm.SetText(password)
	}
	form.
		AddFormItem(svc).
		AddFormItem(login).
		AddFormItem(pwd).
		AddFormItem(confirm).
		AddFormItem(comment).
		AddButton("Ok", func() {
			e.app.QueueUpdateDraw(func() {
				switch {
				case pwd.GetText() == "":
					e.showMessage("Error", "Empty password", screenPwdManage, errorsMessageStyle)
				case pwd.GetText() != confirm.GetText():
					e.showMessage("Error", "Passwords mismatch", screenPwdManage, errorsMessageStyle)
				default:
					p := ezpwd.Password{
						Service:  svc.GetText(),
						Login:    login.GetText(),
						Password: pwd.GetText(),
						Comment:  comment.GetText(),
					}
					if id == -1 {
						pwds = append(pwds, p)
					} else {
						pwds[id] = p
					}
					e.passwordsChan <- pwds
				}
			})
		}).
		AddButton("Cancel", func() {
			e.pages.RemovePage(screenPwdManage)
			e.pages.ShowPage(screenPwds)
		})
	if id == -1 {
		form.SetBorderColor(passwordMgmgFormColors.BorderAdd)
		form.SetTitleColor(passwordMgmgFormColors.TitleAdd)
		form.SetTitle(" Adding new login ")
	} else {
		p := pwds[id]
		form.SetBorderColor(passwordMgmgFormColors.BorderUpdate)
		form.SetTitleColor(passwordMgmgFormColors.TitleUpdate)
		form.SetTitle(fmt.Sprintf(" Updating %s : %s ", p.Service, p.Login))
		svc.SetText(p.Service)
		login.SetText(p.Login)
		pwd.SetText(p.Password)
		confirm.SetText(p.Password)
		comment.SetText(p.Comment)
	}
	form.SetTitleAlign(tview.AlignCenter)
	form.SetBorder(true)
	form.SetCancelFunc(func() {
		e.pages.RemovePage(screenPwdManage)
		e.pages.ShowPage(screenPwds)
	})
	return form
}
