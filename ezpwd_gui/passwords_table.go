package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

func (e *devEzpwd) passwordsTable() {
	table := tview.NewTable().
		SetBorders(true)
	table.SetSelectable(true, false)
	passwordsMsg := tview.NewTextView()
	passwordsMsg.SetBorder(true)
	passwordsMsg.SetDoneFunc(func(tcell.Key) {
		e.pages.RemovePage(screenPwdCopied)
		e.pages.ShowPage(screenPwds)
	})

	table.SetDoneFunc(func(k tcell.Key) {
		switch k {
		case tcell.KeyEsc:
			e.app.Stop()
		}
	})

	go func(t *tview.Table) {
		for pwds := range e.passwordsChan {
			e.app.QueueUpdateDraw(func() {
				t.Clear()
				table.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
					switch key.Rune() {
					case 'a', 'A':
						e.app.QueueUpdateDraw(func() {
							e.pages.AddPage(screenPwdManage,
								modal(e.passwordMgmtForm(-1, pwds), 30, 15),
								true, true,
							)
							e.pages.ShowPage(screenPwdManage)
						})
					case 'u', 'U':
						r, _ := table.GetSelection()
						if r == 0 || len(pwds) == 0 {
							break
						}
						e.app.QueueUpdateDraw(func() {
							e.pages.AddPage(screenPwdManage,
								modal(e.passwordMgmtForm(r-1, pwds), 30, 15),
								true, true,
							)
							e.pages.ShowPage(screenPwdManage)
						})
					case 'd', 'D':
						r, _ := table.GetSelection()
						if r == 0 || r-1 >= len(pwds) {
							break
						}
						e.confirm(fmt.Sprintf(" Remove service '%s : %s'? ", pwds[r-1].Service, pwds[r-1].Login), screenPwds, func() {
							e.app.QueueUpdateDraw(func() {
								e.passwordsChan <- append(pwds[:r-1], pwds[r:]...)
							})
						})
					case 's', 'S':
						{
							err := backup(e.passwordPath)
							switch {
							case err == nil || errors.Is(err, os.ErrNotExist):
								// do nothing
							default:
								e.showMessage("Error", fmt.Sprintf("Can't backup password file: %+v", err), screenPwds)
								break
							}
						}

						var buffer bytes.Buffer

						if err := ezpwd.WritePasswords(pwds, &buffer); err != nil {
							e.showMessage("Error", fmt.Sprintf("Can't backup password file: %+v", err), screenPwds)
							break
						}

						_file, err := os.Create(e.passwordPath)
						if err != nil {
							e.showMessage("Error", fmt.Sprintf("Can't open password file '%s': %+v", e.passwordPath, err), screenPwds)
							break
						}
						defer _file.Close()

						if err := e.crypto.Encrypt(&buffer, _file); err != nil {
							e.showMessage("Error", fmt.Sprintf("Can't encrypt password file '%s': %+v", e.passwordPath, err), screenPwds)
						} else {
							e.showMessage("Success!", fmt.Sprintf("Passwords saved successfully"), screenPwds, func(text *tview.TextView) {
								text.SetBorderColor(tcell.ColorGreen)
								text.SetTitleColor(tcell.ColorGreen)
							})
						}
					}
					return key
				})
				table.SetSelectedFunc(func(row, col int) {
					if row == 0 {
						return
					}
					clipboard.WriteAll(pwds[row-1].Password)
					var content = fmt.Sprintf("Selected password '%s : %s'", pwds[row-1].Service, pwds[row-1].Login)
					passwordsMsg.SetText(content)
					mp := modal(passwordsMsg, len(content)+4, 4)
					e.pages.AddPage(screenPwdCopied, mp, true, true)
					e.pages.ShowPage(screenPwdCopied)
					e.app.SetFocus(passwordsMsg)
				})
				type ColSpec struct {
					name      string
					expansion int
				}
				for i, v := range []ColSpec{{"#", 1}, {"Service", 4}, {"Username", 5}, {"Comment", 10}} {
					table.SetCell(0, i, tview.NewTableCell(v.name).
						SetAlign(tview.AlignCenter).
						SetTextColor(tcell.ColorYellow).
						SetExpansion(v.expansion),
					)
				}
				for i, p := range pwds {
					table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", i+1)).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorRed))
					table.SetCellSimple(i+1, 1, p.Service)
					table.SetCellSimple(i+1, 2, p.Login)
					table.SetCellSimple(i+1, 3, p.Comment)
				}
				table.ScrollToBeginning()
				e.pages.SwitchToPage(screenPwds)
				e.app.SetFocus(table)
				e.app.Draw()
			})
		}
	}(table)
	makeButton := func(txt string) *tview.TextView {
		btn := tview.NewTextView()
		btn.SetBackgroundColor(tcell.ColorBlue)
		return btn.SetText(fmt.Sprintf("[red]%c[white]%s", txt[0], txt[1:])).SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	}
	e.pages.AddPage(screenPwds, tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				AddItem(tview.NewBox(), 0, 1, false).
				AddItem(table, 0, 20, true).
				AddItem(tview.NewBox(), 0, 1, false), 0, 1, true,
		).
		AddItem(
			tview.NewFlex().
				AddItem(tview.NewBox(), 0, 2, false).
				AddItem(makeButton("Add"), 0, 4, true).
				AddItem(tview.NewBox(), 0, 2, false).
				AddItem(makeButton("Update"), 0, 4, true).
				AddItem(tview.NewBox(), 0, 1, false).
				AddItem(makeButton("Delete"), 0, 4, true).
				AddItem(tview.NewBox(), 0, 2, false).
				AddItem(makeButton("Save"), 0, 4, true).
				AddItem(tview.NewBox(), 0, 2, false),
			1, 1, true,
		).SetFullScreen(false), true, false)
}
