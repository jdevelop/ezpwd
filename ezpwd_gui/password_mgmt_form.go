package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

func (e *devEzpwd) passwordMgmtForm(id int, pwds []ezpwd.Password) *tview.Form {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
	svc := tview.NewInputField().SetLabel("Service:").SetFieldWidth(20)
	login := tview.NewInputField().SetLabel("Login:").SetFieldWidth(20)
	pwd := tview.NewInputField().SetLabel("Password:").SetFieldWidth(20)
	confirm := tview.NewInputField().SetLabel("Confirm:").SetFieldWidth(20)
	comment := tview.NewInputField().SetLabel("Comment").SetFieldWidth(20)
	form.
		AddFormItem(svc).
		AddFormItem(login).
		AddFormItem(pwd).
		AddFormItem(confirm).
		AddFormItem(comment).
		AddButton("Ok", func() {
			e.app.QueueUpdateDraw(func() {
				if pwd.GetText() != confirm.GetText() {
					e.showMessage("Passwords mismatch", screenPwdManage)
				} else {
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
		form.SetBorderColor(tcell.ColorBlue)
		form.SetTitle(" Adding new login ")
	} else {
		p := pwds[id]
		form.SetBorderColor(tcell.ColorRed)
		form.SetTitle(fmt.Sprintf(" Updating %s : %s ", p.Service, p.Login)).SetTitleColor(tcell.ColorRed)
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
