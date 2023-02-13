package country

import (
	"errors"

	"github.com/vault-thirteen/auxie/unicode"
)

const (
	ErrCountryCode = "country code error"
	ErrCountryName = "country name error"
)

const (
	CodeUnknown = ""
	NameUnknown = "Unknown"
)

type Country struct {
	code string
	name string
}

func New(code string, name string) (c *Country, err error) {
	if code == CodeUnknown {
		return &Country{
			code: CodeUnknown,
			name: NameUnknown,
		}, nil
	}

	if len(code) != 2 {
		return nil, errors.New(ErrCountryCode)
	}

	if len(name) == 0 {
		return nil, errors.New(ErrCountryName)
	}

	return &Country{
		code: code,
		name: name,
	}, nil
}

func (c *Country) Code() string {
	return c.code
}

func (c *Country) Name() string {
	return c.name
}

func MayBeCountryCode(s string) bool {
	if len(s) != 2 {
		return false
	}
	if !unicode.SymbolIsLatLetter(rune(s[0])) {
		return false
	}
	if !unicode.SymbolIsLatLetter(rune(s[1])) {
		return false
	}

	return true
}
