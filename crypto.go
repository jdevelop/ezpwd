package ezpwd

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/ProtonMail/gopenpgp/constants"
	"github.com/ProtonMail/gopenpgp/crypto"
)

type Crypto struct {
	keyPass []byte
}

func NewCrypto(pwd []byte) (*Crypto, error) {
	return &Crypto{
		keyPass: pwd,
	}, nil
}

type KeyWriter io.Writer

type CryptoInterface interface {
	Encrypt(in io.Reader, out io.Writer) error
	Decrypt(in io.Reader, out io.Writer) error
}

func (cr *Crypto) Encrypt(in io.Reader, out io.Writer) error {
	var key = crypto.NewSymmetricKey(cr.keyPass, constants.ThreeDES)
	content, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	var message = crypto.NewPlainMessage(content)
	encrypted, err := key.Encrypt(message)
	if err != nil {
		return err
	}
	_, err = out.Write(encrypted.GetBinary())
	return err
}

var noSymmetric = errors.New("Symmetric not set")

func (cr *Crypto) Decrypt(in io.Reader, out io.Writer) error {
	var key = crypto.NewSymmetricKey(cr.keyPass, constants.ThreeDES)
	content, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	msg := crypto.NewPGPMessage(content)

	decrypted, err := key.Decrypt(msg)

	if _, err := io.Copy(out, decrypted.NewReader()); err != nil {
		return err
	}
	return nil
}

var _ CryptoInterface = &Crypto{}
