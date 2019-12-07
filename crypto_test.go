package ezpwd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	crypto, err := NewCrypto([]byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	var b, dec bytes.Buffer

	err = crypto.Encrypt(strings.NewReader("test message"), &b)
	require.Nil(t, err)

	err = crypto.Decrypt(&b, &dec)
	require.Nil(t, err)
	require.EqualValues(t, "test message", string(dec.Bytes()))

}

func TestDecryptWrongPass(t *testing.T) {
	crypto, err := NewCrypto([]byte("password"))
	if err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)
	dec := new(bytes.Buffer)

	err = crypto.Encrypt(strings.NewReader("test message"), b)
	require.Nil(t, err)

	crypto.keyPass = []byte("passwor")
	err = crypto.Decrypt(b, dec)
	require.NotNil(t, err)
}
