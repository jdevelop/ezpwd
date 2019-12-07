package ezpwd

import (
	"fmt"
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
		return fmt.Errorf("can't encrypt the message: %w", err)
	}
	if _, err = io.Copy(w, in); err != nil {
		return fmt.Errorf("can't copy encrypted message: %w", err)
	}
	return w.Close()
}

var (
	noSymmetric = fmt.Errorf("Symmetric not set")
	wrongPass   = fmt.Errorf("Wrong password")
)

func (cr *Crypto) Decrypt(in io.Reader, out io.Writer) error {
	read := false
	md, err := openpgp.ReadMessage(in, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if !symmetric {
			return nil, noSymmetric
		}
		if read {
			return nil, wrongPass
		}
		read = true
		return cr.keyPass, nil
	}, nil)
	if err != nil {
		return fmt.Errorf("can't decrypt message : %w", err)
	}
	if _, err := io.Copy(out, md.UnverifiedBody); err != nil {
		return fmt.Errorf("can't transfer decrypted message : %w", err)
	}
	return nil
}

var _ CryptoInterface = &Crypto{}
