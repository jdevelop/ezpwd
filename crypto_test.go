package ezpwd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	crypto, err := NewCrypto([]byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)
	dec := new(bytes.Buffer)

	err = crypto.Encrypt(strings.NewReader("test message"), b)
	if err != nil {
		t.Fatal(err)
	}

	err = crypto.Decrypt(b, dec)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, "test message", string(dec.Bytes()))

}
