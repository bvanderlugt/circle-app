package main

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxNameLen = 50

var (
	errNameRequired    = errors.New("name is required")
	errNameTooLong     = errors.New("name must be 50 characters or fewer")
	errNameInvalidChar = errors.New("name contains invalid characters")
)

func validateName(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errNameRequired
	}
	if utf8.RuneCountInString(s) > maxNameLen {
		return "", errNameTooLong
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return "", errNameInvalidChar
		}
	}
	return s, nil
}
