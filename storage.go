package ezpwd

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var splitter = regexp.MustCompile("\\s*/\\s*")

type Password struct {
	Service  string
	Login    string
	Password string
	Comment  string
}

func ReadPasswords(src io.Reader) ([]Password, error) {
	passwords := make([]Password, 0, 50)
	rdr := bufio.NewScanner(src)
	for rdr.Scan() {
		line := rdr.Text()
		parts := splitter.Split(line, 4)
		pLen := len(parts)
		if pLen < 3 || pLen > 4 {
			continue
		}
		p := Password{
			Service:  parts[0],
			Login:    parts[1],
			Password: parts[2],
		}
		if pLen == 4 {
			p.Comment = parts[3]
		}
		passwords = append(passwords, p)
	}
	return passwords, nil
}

func WritePasswords(passwords []Password, writer io.Writer) error {
	wrtr := bufio.NewWriter(writer)
	for _, pwd := range passwords {
		if pwd.Comment != "" {
			if _, err := wrtr.WriteString(fmt.Sprintf("%s / %s / %s / %s \n", pwd.Service, pwd.Login, pwd.Password, pwd.Comment)); err != nil {
				return err
			}
		} else {
			if _, err := wrtr.WriteString(fmt.Sprintf("%s / %s / %s \n", pwd.Service, pwd.Login, pwd.Password)); err != nil {
				return err
			}
		}
	}
	return wrtr.Flush()
}
