package main

import (
	"errors"
	"regexp"
	"strings"
)

var nameMatch = regexp.MustCompile(`\A((@[^\s\/~'!\(\)\*]+?)[\/])?([^_.][^\s\/~'!\(\)\*]+)\z`)

func IsValidName(name string) error {
	if strings.TrimSpace(name) != name {
		return errors.New("username contains space")
	}
	if len(name) == 0 || len(name) > Conf.Common.MaxUsernameLen {
		return errors.New("username invalid len")
	}
	if !nameMatch.MatchString(name) {
		return errors.New("username invalid pattern")
	}
	return nil
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+-/=?^_`{|}~]*@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsValidEmail(email string) error {
	if len(email) == 0 {
		return errors.New("invalid email len")
	}

	if email[0] == '-' {
		return errors.New("invalid email")
	}

	n := strings.LastIndex(email, "@")
	if n <= 0 {
		return errors.New("invalid email address")
	}

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email pattern")
	}

	return nil
}

func IsValidPasswd(passwd string) error {
	if len(passwd) < Conf.Common.MinPasswordLength {
		return errors.New("password too short")
	}
	if !IsComplexEnough(passwd) {
		return errors.New("password too simple")
	}

	return nil
}
