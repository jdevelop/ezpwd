package ezpwd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expected = []Password{
	{
		Service:  "service",
		Login:    "username",
		Password: "password",
		Comment:  "",
	},
	{
		Service:  "another service",
		Login:    "user",
		Password: "pwd",
		Comment:  "comment",
	},
}

func TestPasswordRead(t *testing.T) {
	rdr := strings.NewReader(`service / username / password
another service / user / pwd / comment
no service
#comment`)
	pwds, err := ReadPasswords(rdr)
	assert.Nil(t, err)
	assert.EqualValues(t, expected[:], pwds)
}

func TestPasswordWrite(t *testing.T) {
	b := new(bytes.Buffer)
	err := WritePasswords(expected, b)

	assert.Nil(t, err)

	assert.Equal(t, `service / username / password 
another service / user / pwd / comment 
`, b.String())
}
