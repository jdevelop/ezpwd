package ezpwd

import (
	"io"

	"golang.org/x/crypto/openpgp"
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
	w, err := openpgp.SymmetricallyEncrypt(out, cr.keyPass, nil, nil)
	if err != nil {
		return err
	}
	if _, err = io.Copy(w, in); err != nil {
		return err
	}
	return w.Close()
}

func (cr *Crypto) Decrypt(in io.Reader, out io.Writer) error {
	md, err := openpgp.ReadMessage(in, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return cr.keyPass, nil
	}, nil)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, md.UnverifiedBody); err != nil {
		return err
	}
	return nil
}

var _ CryptoInterface = &Crypto{}
