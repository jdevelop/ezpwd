package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

const (
	screenPwd       = "password"
	screenPwds      = "passwords"
	screenError     = "errors"
	screenPwdCopied = "copied"
	screenPwdManage = "manage"
	screenConfirm   = "confirm"
)

var (
	passFile = flag.String("passfile", "private/test-pass.enc", "Password file")
)

type devEzpwd struct {
	app           *tview.Application
	screen        tcell.Screen
	passwordPath  string
	passwordsChan chan []ezpwd.Password
	pages         *tview.Pages
}

func NewEzpwd(passwordsPath string) (*devEzpwd, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}
	app := tview.NewApplication()
	app.SetScreen(screen)

	var instance = devEzpwd{
		app:          app,
		screen:       screen,
		passwordPath: passwordsPath,
		pages:        tview.NewPages(),
	}

	return &instance, nil
}

func (e *devEzpwd) Run() error {
	e.passwordsChan = make(chan []ezpwd.Password)
	form := e.passwordForm(func(pwd string) {
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
				e.passwordsChan <- pwds
			}
		case errors.Is(err, os.ErrNotExist):
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

	})
	tableContainer := e.passwordsTable()
	e.pages.
		AddPage(screenPwd, modal(form, 40, 8), true, true).
		AddPage(screenPwds, tableContainer, true, false)
	e.app.SetRoot(e.pages, true).SetFocus(form)
	return e.app.Run()
}

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal("can't retrieve current user", err)
	}
	flag.Parse()
	encPath := filepath.Join(u.HomeDir, *passFile)

	p, err := NewEzpwd(encPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

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
						if r == 0 {
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

func (e *devEzpwd) passwordForm(onComplete func(string)) tview.Primitive {
	form := tview.NewForm().SetButtonsAlign(tview.AlignCenter)
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
