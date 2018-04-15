package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"os/user"

	"github.com/jdevelop/ezpwd"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	add := flag.Bool("add", false, "Add new password")
	passFile := flag.String("passfile", "private/test-pass.enc", "Password file")

	flag.Parse()

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	encPath := filepath.Join(u.HomeDir, *passFile)
	fmt.Print("Storage Password :/> ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()

	crypto, err := ezpwd.NewCrypto(bytePassword)
	if err != nil {
		log.Fatal(err)
	}

	if *add {
		if err := addFunc(crypto, encPath); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := listFunc(crypto, encPath); err != nil {
			log.Fatal(err)
		}
	}

}

func listPasswords(cr ezpwd.CryptoInterface, file string) ([]ezpwd.Password, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(file); err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	if err := cr.Decrypt(f, buffer); err != nil {
		return nil, err
	}

	return ezpwd.ReadPasswords(buffer)
}

func listFunc(cr ezpwd.CryptoInterface, file string) error {
	pwds, err := listPasswords(cr, file)
	if err != nil {
		return err
	}
	for _, pwd := range pwds {
		if pwd.Comment != "" {
			fmt.Printf("%s / %s / %s / %s\n", pwd.Service, pwd.Login, pwd.Password, pwd.Comment)
		} else {
			fmt.Printf("%s / %s / %s\n", pwd.Service, pwd.Login, pwd.Password)
		}
	}

	return nil
}

func addFunc(cr ezpwd.CryptoInterface, file string) error {
	var (
		_file *os.File
		pwds  []ezpwd.Password
	)
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			if _file, err = os.Create(file); err != nil {
				return err
			}
		} else {
			return err
		}
		pwds = make([]ezpwd.Password, 0, 1)
	} else {
		if _file, err = os.Open(file); err != nil {
			return err
		}
		if pwds, err = listPasswords(cr, file); err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	var service, email, password, comment string

	fmt.Print("Service :/> ")
	if scanner.Scan() {
		service = scanner.Text()
	}

	fmt.Print("Username/email :/> ")
	if scanner.Scan() {
		email = scanner.Text()
	}
	for {
		fmt.Print("Enter Password :/> ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Print("\nConfirm Password :/> ")
		confirmPassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		if string(bytePassword) == string(confirmPassword) {
			password = string(bytePassword)
			break
		} else {
			fmt.Println()
		}
	}

	fmt.Print("\nComment :/> ")
	if scanner.Scan() {
		comment = scanner.Text()
	}

	pwds = append(pwds, ezpwd.Password{
		Service:  service,
		Login:    email,
		Password: password,
		Comment:  comment,
	})

	if err := backup(file); err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	err := ezpwd.WritePasswords(pwds, buffer)

	if err != nil {
		return err
	}

	if err = _file.Close(); err != nil {
		return err
	}

	_file, err = os.OpenFile(file, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	return cr.Encrypt(buffer, _file)

}

const layout = "020106"

func backup(file string) error {
	current := time.Now()
	newName := fmt.Sprintf("%s.%s-%d", file, current.Format(layout), current.Unix())
	src, err := os.Open(file)
	if err != nil {
		return err
	}
	f, err := os.Create(newName)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, src)
	return err
}
