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
	"strconv"
	"syscall"
	"time"

	"os/user"

	"github.com/atotto/clipboard"
	"github.com/jdevelop/ezpwd"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	add      = flag.Bool("add", false, "Add new password")
	list     = flag.Bool("list", false, "List all passwords")
	passFile = flag.String("passfile", "private/pass.enc", "Password file")
	upd      = flag.Bool("update", false, "Update password")
)

func main() {

	flag.Parse()

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	encPath := filepath.Join(u.HomeDir, *passFile)
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

	switch {
	case *add:
		if err := addFunc(crypto, encPath); err != nil {
			log.Fatal(err)
		}
	case *list:
		rdr := ezpwd.NewShownPWDs(os.Stdout)
		if err := listFunc(crypto, rdr, encPath); err != nil {
			log.Fatal(err)
		}
	case *upd:
		rdr := ezpwd.NewHiddenPWDs(os.Stdout)
		if err := updateFunc(crypto, rdr, encPath); err != nil {
			log.Fatal(err)
		}
	default:
		rdr := ezpwd.NewHiddenPWDs(os.Stdout)
		if err := listAndCopyFunc(crypto, rdr, encPath); err != nil {
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

func listAndCopyFunc(cr ezpwd.CryptoInterface, rdr ezpwd.Renderer, file string) error {
	pwds, err := listPasswords(cr, file)
	if err != nil {
		return err
	}
	rdr.RenderPasswords(pwds)
	fmt.Printf("Choose password ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	selection, err := strconv.ParseInt(s.Text(), 10, 64)
	if err != nil {
		return err
	}
	if selection >= 0 && int(selection) < len(pwds) {
		return clipboard.WriteAll(pwds[selection].Password)
	}
	return nil
}

func listFunc(cr ezpwd.CryptoInterface, rdr ezpwd.Renderer, file string) error {
	pwds, err := listPasswords(cr, file)
	if err != nil {
		return err
	}
	rdr.RenderPasswords(pwds)
	return nil
}

func readInput(pwd *ezpwd.Password) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Service :/> ")
	if scanner.Scan() {
		if s := scanner.Text(); s != "" {
			pwd.Service = s
		}
	}

	fmt.Print("Username/email :/> ")
	if scanner.Scan() {
		if s := scanner.Text(); s != "" {
			pwd.Login = s
		}
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
			if s := string(bytePassword); s != "" {
				pwd.Password = s
			}
			break
		} else {
			fmt.Println()
		}
	}

	fmt.Print("\nComment :/> ")
	if scanner.Scan() {
		if s := scanner.Text(); s != "" {
			pwd.Comment = s
		}
	}
	return nil
}

func readPasswords(cr ezpwd.CryptoInterface, file string) ([]ezpwd.Password, *os.File, error) {
	var (
		_file *os.File
		pwds  []ezpwd.Password
	)
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			if _file, err = os.Create(file); err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
		pwds = make([]ezpwd.Password, 0, 1)
	} else {
		if _file, err = os.Open(file); err != nil {
			return nil, nil, err
		}
		if pwds, err = listPasswords(cr, file); err != nil {
			return nil, nil, err
		}
	}
	return pwds, _file, nil
}

func updateFunc(cr ezpwd.CryptoInterface, render ezpwd.Renderer, file string) error {
	pwds, _file, err := readPasswords(cr, file)
	if err != nil {
		return err
	}
	render.RenderPasswords(pwds)
	scanner := bufio.NewScanner(os.Stdin)
	var (
		idx int
	)
	for {
		fmt.Print("Please choose the entry you'd like to change: ")
		scanner.Scan()
		numStr := scanner.Text()
		_idx, err := strconv.ParseInt(numStr, 10, 32)
		if err != nil {
			return err
		}
		idx = int(_idx)
		if l := len(pwds); l <= idx {
			fmt.Printf("Enter numbers between 0 and %d\n", l)
		} else {
			break
		}
	}
	err = readInput(&pwds[idx])
	if err != nil {
		return err
	}
	if err := backup(file); err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err := ezpwd.WritePasswords(pwds, buffer); err != nil {
		return err
	}
	if err := _file.Close(); err != nil {
		return err
	}
	_file, err = os.OpenFile(file, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	return cr.Encrypt(buffer, _file)
}

func addFunc(cr ezpwd.CryptoInterface, file string) error {

	pwds, _file, err := readPasswords(cr, file)
	if err != nil {
		return err
	}

	var pwd ezpwd.Password

	err = readInput(&pwd)

	if err != nil {
		return err
	}

	pwds = append(pwds, pwd)

	if err := backup(file); err != nil {
		return err
	}

	buffer := new(bytes.Buffer)

	if err := ezpwd.WritePasswords(pwds, buffer); err != nil {
		return err
	}

	if err := _file.Close(); err != nil {
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
