package country

import (
	"github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/unicode"
)

const (
	Err_CountryCode = "country code error"
	Err_CountryName = "country name error"
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
		return nil, errors.ErrorFromText(Err_CountryCode)
	}

	if len(name) == 0 {
		return nil, errors.ErrorFromText(Err_CountryName)
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
