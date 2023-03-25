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
		return errors.New("invalid_email_len")
	}

	if email[0] == '-' {
		return errors.New("invalid_email")
	}

	n := strings.LastIndex(email, "@")
	if n <= 0 {
		return errors.New("invalid_email_address")
	}

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid_email_pattern")
	}

	return nil
}

func IsValidPasswd(passwd string) error {
	if len(passwd) < Conf.Common.MinPasswordLength {
		return errors.New("password_too_short")
	}
	if !IsComplexEnough(passwd) {
		return errors.New("password_too_simple")
	}

	return nil
}
