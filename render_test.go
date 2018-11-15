package ezpwd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pwds = []Password{
	{
		Service:  "Gmail",
		Login:    "123@123.com",
		Password: "password",
	},
	{
		Service:  "Yahoo",
		Login:    "123@yahoo.com",
		Password: "nopassword",
	},
}

func TestRenderNoPassword(t *testing.T) {
	expected := `+---+---------+---------------+---------+
| # | SERVICE |     LOGIN     | COMMENT |
+---+---------+---------------+---------+
| 0 | Gmail   | 123@123.com   |         |
| 1 | Yahoo   | 123@yahoo.com |         |
+---+---------+---------------+---------+
`
	w := bytes.Buffer{}
	pwdWriter := NewHiddenPWDs(&w)
	pwdWriter.RenderPasswords(pwds)
	content := string(w.Bytes())
	assert.Equal(t, expected, content)
}

func TestRenderAllPassword(t *testing.T) {
	expected := `+---+---------+---------------+------------+---------+
| # | SERVICE |     LOGIN     |  PASSWORD  | COMMENT |
+---+---------+---------------+------------+---------+
| 0 | Gmail   | 123@123.com   | password   |         |
| 1 | Yahoo   | 123@yahoo.com | nopassword |         |
+---+---------+---------------+------------+---------+
`
	w := bytes.Buffer{}
	pwdWriter := NewShownPWDs(&w)
	pwdWriter.RenderPasswords(pwds)
	content := string(w.Bytes())
	assert.Equal(t, expected, content)
}
