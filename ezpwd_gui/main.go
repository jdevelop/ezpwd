package main

import (
	"flag"
	"log"
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
	crypto        *ezpwd.Crypto
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
	form := e.passwordForm()
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
