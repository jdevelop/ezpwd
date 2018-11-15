package ezpwd

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

type Renderer interface {
	RenderPasswords([]Password)
}

type hidePasswords struct {
	w io.Writer
}

type showPasswords struct {
	w io.Writer
}

func NewHiddenPWDs(w io.Writer) *hidePasswords {
	return &hidePasswords{w: w}
}

func NewShownPWDs(w io.Writer) *showPasswords {
	return &showPasswords{w: w}
}

func (hp *hidePasswords) RenderPasswords(pwds []Password) {
	table := tablewriter.NewWriter(hp.w)
	for i, pwd := range pwds {
		table.Append([]string{fmt.Sprintf("%d", i), pwd.Service, pwd.Login, pwd.Comment})
	}
	table.SetHeader([]string{"#", "Service", "Login", "Comment"})
	table.Render()
}

func (sp *showPasswords) RenderPasswords(pwds []Password) {
	table := tablewriter.NewWriter(sp.w)
	for i, pwd := range pwds {
		table.Append([]string{fmt.Sprintf("%d", i), pwd.Service, pwd.Login, pwd.Password, pwd.Comment})
	}
	table.SetHeader([]string{"#", "Service", "Login", "Password", "Comment"})
	table.Render()
}

var _ Renderer = &showPasswords{}
var _ Renderer = &hidePasswords{}
