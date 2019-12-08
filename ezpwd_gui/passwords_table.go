package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (e *devEzpwd) passwordsTable() tview.Primitive {
	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 0)
	for i, v := range []string{"#", "Service", "Username", "Comment"} {
		table.SetCell(0, i, tview.NewTableCell(v).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorYellow))
	}
	table.SetSelectable(true, false)
	passwordsMsg := tview.NewTextView()
	passwordsMsg.SetBorder(true)
	pages := tview.NewPages()
	passwordsMsg.SetDoneFunc(func(tcell.Key) {
		pages.RemovePage(screenPwdCopied)
		pages.ShowPage(screenPwds)
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
						/*
							if err := backup(file); err != nil {
								return err
							}

							buffer := new(bytes.Buffer)

							if err := ezpwd.WritePasswords(pwds, buffer); err != nil {
								return err
							}

							if err := _file.Close(); err != nil {
								return err
							}

							_file, err = os.OpenFile(file, os.O_WRONLY, os.ModeAppend)
							if err != nil {
								return err
							}

							return cr.Encrypt(buffer, _file)
						*/
					}
					return key
				})
				table.SetSelectedFunc(func(row, col int) {
					if row == 0 {
						return
					}
					var content = fmt.Sprintf("Selected password '%s : %s'", pwds[row-1].Service, pwds[row-1].Login)
					passwordsMsg.SetText(content)
					mp := modal(passwordsMsg, len(content)+4, 4)
					pages.AddPage(screenPwdCopied, mp, true, true)
					pages.ShowPage(screenPwdCopied)
					e.app.SetFocus(passwordsMsg)
				})
				for i, v := range []string{"#", "Service", "Username", "Comment"} {
					table.SetCell(0, i, tview.NewTableCell(v).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorYellow))
				}
				for i, p := range pwds {
					table.SetCellSimple(i+1, 0, fmt.Sprintf("%d", i+1))
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
	pages.AddPage(screenPwds, tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				AddItem(tview.NewBox(), 0, 2, false).
				AddItem(table, 0, 2, true).
				AddItem(tview.NewBox(), 0, 2, false), 0, 1, true,
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
		).SetFullScreen(true), true, true)
	return pages
}
