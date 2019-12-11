package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/jdevelop/ezpwd"
	"github.com/rivo/tview"
)

const (
	screenPwd       = "password"
	screenPwds      = "passwordsTable"
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
	var (
		form   *tview.Form
		height int
	)
	_, err := os.Stat(e.passwordPath)
	switch {
	case err == nil:
		form, height = e.passwordForm(), 7
	case errors.Is(err, os.ErrNotExist):
		form, height = e.initForm(), 9
	default:
		log.Fatalf("can't use storage at path %s: %+v", e.passwordPath, err)
	}
	e.passwordsTable()
	e.pages.
		AddPage(screenPwd, modal(form, 40, height, func(p *tview.Box) {
			p.SetBackgroundColor(globalScreenColors.Background)
			p.SetBorderColor(globalScreenColors.Border)
			p.SetTitleColor(globalScreenColors.Title)
		}), true, true)
	frame := tview.NewFrame(e.pages)
	frame.SetBorder(true)
	frame.SetBackgroundColor(globalScreenColors.Background)
	frame.SetBorderColor(globalScreenColors.Border)
	frame.SetTitleColor(globalScreenColors.Title)
	frame.AddText("`Esc` to exit dialogs without saving", false, tview.AlignRight, globalScreenColors.HelpText)
	frame.AddText("`Ctrl-C` to quit application", false, tview.AlignRight, globalScreenColors.HelpText)
	frame.SetTitle("Storage " + e.passwordPath)
	e.app.SetRoot(frame, true).SetFocus(form)
	return e.app.Run()
}

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal("can't retrieve current user", err)
	}
	flag.Parse()
	var encPath string
	if !strings.HasPrefix(*passFile, "/") {
		encPath = filepath.Join(u.HomeDir, *passFile)
	} else {
		encPath = *passFile
	}

	p, err := NewEzpwd(encPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
