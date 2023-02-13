package ipa

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	ErrIPAddressV4Syntax = "syntax error in IPv4 address"
	ErrOverflow          = "overflow"
)

type IPAddressV4 uint32

func NewFromBytes(ba, bb, bc, bd byte) IPAddressV4 {
	return (IPAddressV4(ba) << 24) + (IPAddressV4(bb) << 16) +
		(IPAddressV4(bc) << 8) + IPAddressV4(bd)
}

func NewFromString(str string) (addr IPAddressV4, err error) {
	parts := strings.Split(str, ".")
	if len(parts) != 4 {
		return addr, errors.New(ErrIPAddressV4Syntax)
	}

	var bytes = make([]byte, 0, 4)
	var tmp uint64
	for _, part := range parts {
		tmp, err = strconv.ParseUint(part, 10, 64)
		if err != nil {
			return addr, err
		}
		if tmp > math.MaxUint8 {
			return addr, errors.New(ErrOverflow)
		}

		bytes = append(bytes, byte(tmp))
	}

	return NewFromBytes(bytes[0], bytes[1], bytes[2], bytes[3]), nil
}

func NewFromUintString(uintStr string) (addr IPAddressV4, err error) {
	var u uint64
	u, err = strconv.ParseUint(uintStr, 10, 64)
	if err != nil {
		return 0, err
	}

	if u > math.MaxUint32 {
		return addr, errors.New(ErrOverflow)
	}

	return IPAddressV4(u), nil
}
