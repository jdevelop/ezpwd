package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

func (e *devEzpwd) passwordsTable() {
	table := tview.NewTable().
		SetBorders(true)
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetBackgroundColor(DefaultColorSchema.PasswordsTableColors.Background)
	table.SetTitleColor(DefaultColorSchema.PasswordsTableColors.Title)
	table.SetBorderColor(DefaultColorSchema.PasswordsTableColors.Border)
	table.SetBordersColor(DefaultColorSchema.PasswordsTableColors.Border)
	table.SetSelectedStyle(DefaultColorSchema.PasswordsTableColors.Selection, DefaultColorSchema.PasswordsTableColors.SelectionBackground, 0)
	filterBox := tview.NewInputField().SetLabel("Filter: ")
	filterBox.SetBackgroundColor(DefaultColorSchema.PasswordsTableColors.Background)
	filterBox.SetFieldTextColor(DefaultColorSchema.PasswordsTableColors.FieldText)
	filterBox.SetFieldBackgroundColor(DefaultColorSchema.PasswordsTableColors.FieldBackground)
	filterBox.SetLabelColor(DefaultColorSchema.PasswordsTableColors.Label)
	passwordsMsg := tview.NewTextView()
	passwordsMsg.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	passwordsMsg.SetBackgroundColor(DefaultColorSchema.PasswordsTableColors.CopiedBackground)
	passwordsMsg.SetTextColor(DefaultColorSchema.PasswordsTableColors.CopiedText)
	passwordsMsg.SetTitleColor(DefaultColorSchema.PasswordsTableColors.CopiedTitle)
	passwordsMsg.SetBorderColor(DefaultColorSchema.PasswordsTableColors.CopiedBorder)
	passwordsMsg.SetDoneFunc(func(tcell.Key) {
		e.pages.RemovePage(screenPwdCopied)
		e.pages.ShowPage(screenPwds)
	})

	go func(t *tview.Table) {
		var (
			currentPasswords *[]ezpwd.Password
			drawTable        func(string)
			mappings         []int
		)

		t.SetDoneFunc(func(k tcell.Key) {
			switch k {
			case tcell.KeyEsc:
				if filterBox.GetText() != "" {
					e.app.QueueUpdateDraw(func() {
						filterBox.SetText("")
						drawTable("")
						e.app.SetFocus(table)
						e.app.Draw()
					})
				} else {
					e.confirm(" Are you sure you want to quit? ", screenPwds, e.app.Stop)
				}
			}
		})

		filterBox.SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				e.app.QueueUpdateDraw(func() {
					drawTable(filterBox.GetText())
					e.app.SetFocus(table)
					e.app.Draw()
				})
			case tcell.KeyEsc:
				filterBox.SetText("")
				e.app.QueueUpdateDraw(func() {
					drawTable("")
					e.app.SetFocus(table)
					e.app.Draw()
				})
			}
		})
		drawTable = func(filter string) {
			t.Clear()
			mappings = make([]int, 0)
			table.SetSelectedFunc(func(row, col int) {
				if row == 0 {
					return
				}
				row = mappings[row-1]
				clipboard.WriteAll((*currentPasswords)[row].Password)
				var content = fmt.Sprintf(" Password copied to clipboard '%s : %s' ", (*currentPasswords)[row].Service, (*currentPasswords)[row].Login)
				passwordsMsg.SetText(content)
				mp := modal(passwordsMsg, len(content)+2, 3)
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
					SetTextColor(DefaultColorSchema.PasswordsTableColors.Header).
					SetExpansion(v.expansion),
				)
			}
			i := 1
			equals := func(src, substr string) bool {
				return strings.Contains(strings.ToUpper(src), strings.ToUpper(substr))
			}
			for idx, p := range *currentPasswords {
				if filter != "" && !(equals(p.Service, filter) || equals(p.Comment, filter) || equals(p.Login, filter)) {
					continue
				}
				mappings = append(mappings, idx)
				table.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("%d", i)).SetAlign(tview.AlignCenter))
				table.SetCellSimple(i, 1, p.Service)
				table.SetCellSimple(i, 2, p.Login)
				table.SetCellSimple(i, 3, p.Comment)
				i += 1
			}
			table.ScrollToBeginning()

		}
		dialogsStyle := func(b *tview.Box) {
			b.SetBackgroundColor(DefaultColorSchema.PasswordsTableColors.Background)
		}
		table.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
			switch key.Rune() {
			case 'a', 'A':
				e.app.QueueUpdateDraw(func() {
					e.pages.AddPage(screenPwdManage,
						modal(e.passwordMgmtForm(-1, *currentPasswords), 50, 15, dialogsStyle),
						true, true,
					)
					e.pages.ShowPage(screenPwdManage)
				})
			case 'u', 'U':
				r, _ := table.GetSelection()
				if r == 0 || len(*currentPasswords) == 0 {
					break
				}
				e.app.QueueUpdateDraw(func() {
					e.pages.AddPage(screenPwdManage,
						modal(e.passwordMgmtForm(r-1, *currentPasswords), 50, 15, dialogsStyle),
						true, true,
					)
					e.pages.ShowPage(screenPwdManage)
				})
			case 'd', 'D':
				r, _ := table.GetSelection()
				if r == 0 || r-1 >= len(*currentPasswords) {
					break
				}
				r = mappings[r-1]
				e.confirm(fmt.Sprintf(" Remove service '%s : %s'? ", (*currentPasswords)[r].Service, (*currentPasswords)[r].Login), screenPwds, func() {
					e.app.QueueUpdateDraw(func() {
						*currentPasswords = append((*currentPasswords)[:r], (*currentPasswords)[r+1:]...)
						drawTable(filterBox.GetText())
					})
				})
			case 'f', 'F':
				e.app.SetFocus(filterBox)
			case 's', 'S':
				{
					err := backup(e.passwordPath)
					switch {
					case err == nil || errors.Is(err, os.ErrNotExist):
						// do nothing
					default:
						e.showMessage("Error", fmt.Sprintf("Can't backup password file: %+v", err), screenPwds, errorsMessageStyle)
						break
					}
				}
				var buffer bytes.Buffer
				if err := ezpwd.WritePasswords(*currentPasswords, &buffer); err != nil {
					e.showMessage("Error", fmt.Sprintf("Can't backup password file: %+v", err), screenPwds, errorsMessageStyle)
					break
				}
				_file, err := os.Create(e.passwordPath)
				if err != nil {
					e.showMessage("Error", fmt.Sprintf("Can't open password file '%s': %+v", e.passwordPath, err), screenPwds, errorsMessageStyle)
					break
				}
				defer _file.Close()
				if err := e.crypto.Encrypt(&buffer, _file); err != nil {
					e.showMessage("Error", fmt.Sprintf("Can't encrypt password file '%s': %+v", e.passwordPath, err), screenPwds, errorsMessageStyle)
				} else {
					e.showMessage("Success!", fmt.Sprintf("Passwords saved successfully"), screenPwds, successMessageStyle)
				}
			}
			return key
		})

		for pwds := range e.passwordsChan {
			currentPasswords = &pwds
			drawTable("")
			e.app.QueueUpdateDraw(func() {
				e.pages.SwitchToPage(screenPwds)
				e.app.SetFocus(table)
				e.app.Draw()
			})
		}
	}(table)
	makeButton := func(txt string) *tview.TextView {
		btn := tview.NewTextView()
		btn.SetBackgroundColor(DefaultColorSchema.PasswordsTableColors.ButtonBackground)
		btn.SetTextColor(DefaultColorSchema.PasswordsTableColors.ButtonText)
		return btn.SetText(fmt.Sprintf("[%s]%c[%s]%s", color2Hex(DefaultColorSchema.PasswordsTableColors.ButtonAccent), txt[0],
			color2Hex(DefaultColorSchema.PasswordsTableColors.ButtonText), txt[1:])).
			SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	}
	flexColors := func(flex *tview.Flex) *tview.Flex {
		flex.SetBackgroundColor(DefaultColorSchema.PasswordsTableColors.Background)
		flex.SetTitleColor(DefaultColorSchema.PasswordsTableColors.Title)
		flex.SetBorderColor(DefaultColorSchema.PasswordsTableColors.Background)
		return flex
	}
	spacerBox := func() *tview.Box {
		return tview.NewBox().SetBackgroundColor(DefaultColorSchema.PasswordMgmtColors.Background)
	}
	e.pages.AddPage(screenPwds, flexColors(tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(filterBox, 2, 0, false).
		AddItem(
			flexColors(tview.NewFlex().
				AddItem(spacerBox(), 0, 1, false).
				AddItem(table, 0, 20, true).
				AddItem(spacerBox(), 0, 1, false)), 0, 1, true,
		)).
		AddItem(
			flexColors(tview.NewFlex().
				AddItem(spacerBox(), 0, 2, false).
				AddItem(makeButton("Filter"), 0, 4, true).
				AddItem(spacerBox(), 0, 2, false).
				AddItem(makeButton("Add"), 0, 4, true).
				AddItem(spacerBox(), 0, 2, false).
				AddItem(makeButton("Update"), 0, 4, true).
				AddItem(spacerBox(), 0, 1, false).
				AddItem(makeButton("Delete"), 0, 4, true).
				AddItem(spacerBox(), 0, 2, false).
				AddItem(makeButton("Save"), 0, 4, true).
				AddItem(spacerBox(), 0, 2, false)),
			1, 1, true,
		).SetFullScreen(false), true, false)
}
