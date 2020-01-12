package main

import (
	"errors"
	"flag"
	easyjson "github.com/mailru/easyjson"
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
	passFile   = flag.String("passfile", "private/pass.enc", "Password file")
	schemaFile = flag.String("schema", "", "Color schema file")
	dumpSchema = flag.Bool("dump-schema", false, "Print current schema and exit")
	darkSchema = flag.Bool("dark", false, "Dark colors")
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
			p.SetBackgroundColor(DefaultColorSchema.GlobalScreenColors.Background)
			p.SetBorderColor(DefaultColorSchema.GlobalScreenColors.Border)
			p.SetTitleColor(DefaultColorSchema.GlobalScreenColors.Title)
		}), true, true)
	frame := tview.NewFrame(e.pages)
	frame.SetBorder(true)
	frame.SetBackgroundColor(DefaultColorSchema.GlobalScreenColors.Background)
	frame.SetBorderColor(DefaultColorSchema.GlobalScreenColors.Border)
	frame.SetTitleColor(DefaultColorSchema.GlobalScreenColors.Title)
	frame.AddText("`Esc` to exit dialogs without saving", false, tview.AlignRight, DefaultColorSchema.GlobalScreenColors.HelpText)
	frame.AddText("`Ctrl-C` to quit application", false, tview.AlignRight, DefaultColorSchema.GlobalScreenColors.HelpText)
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

	dir := filepath.Dir(encPath)

	switch d, err := os.Stat(dir); err {
	case nil:
		if !d.IsDir() {
			log.Fatalf("Can't use folder '%s'", dir)
		}
	case os.ErrNotExist:
		if err := os.MkdirAll(dir, 0700); err != nil {
			log.Fatalf("Can't create folder '%s'", dir)
		}
	default:
		log.Fatalf("Fatal error, aborting: %+v", err)
	}

	switch {
	case *schemaFile != "":
		f, err := os.Open(*schemaFile)
		if err != nil {
			log.Fatalf("Can't open schema file %s: %+v", *schemaFile, err)
		}
		if err := easyjson.UnmarshalFromReader(f, &DefaultColorSchema); err != nil {
			log.Fatalf("Can't read JSON from %s: %+v", *schemaFile, err)
		}
	case *darkSchema:
		DefaultColorSchema = DarkColorSchema
	default:
		DefaultColorSchema = LightColorSchema
	}

	if *dumpSchema {
		if _, err := easyjson.MarshalToWriter(&DefaultColorSchema, os.Stdout); err != nil {
			log.Fatalf("Can't encode JSON: %+v", err)
		}
		os.Exit(0)
	}

	p, err := NewEzpwd(encPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
