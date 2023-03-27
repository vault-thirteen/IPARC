package ipar

import (
	"errors"
	"math"

	"github.com/vault-thirteen/IPARC/ipad"
	"github.com/vault-thirteen/IPARC/ipad/country"
	"github.com/vault-thirteen/auxie/IPA"
)

const (
	ErrStartIsReversed = "start is reversed"
	ErrDataIsNotSet    = "data is not set"
)

type IPAddressV4Range struct {
	start  float64
	end    float64
	middle float64
	radius float64
	data   *ipad.IPAddressData
}

func New(
	start ipa.IPAddressV4,
	end ipa.IPAddressV4,
	data *ipad.IPAddressData,
) (r *IPAddressV4Range, err error) {
	if end < start {
		return nil, errors.New(ErrStartIsReversed)
	}
	if data == nil {
		return nil, errors.New(ErrDataIsNotSet)
	}

	r = &IPAddressV4Range{
		start: float64(start),
		end:   float64(end),
		data:  data,
	}

	r.middle = (r.start + r.end) / 2
	r.radius = r.middle - r.start

	return r, nil
}

func (r *IPAddressV4Range) Contains(value float64) (contains bool) {
	return (r.start <= value) && (value <= r.end)
}

func (r *IPAddressV4Range) GetStart() (mid float64) {
	return r.start
}

func (r *IPAddressV4Range) GetEnd() (mid float64) {
	return r.end
}

func (r *IPAddressV4Range) GetMiddle() (mid float64) {
	return r.middle
}

func (r *IPAddressV4Range) GetRadius() (rad float64) {
	return r.radius
}

func (r *IPAddressV4Range) HasIntersectionWith(that *IPAddressV4Range) (intersects bool) {
	return math.Abs(r.middle-that.middle) <= (r.radius + that.radius)
}

func IsSequence(r1 *IPAddressV4Range, r2 *IPAddressV4Range) (isSequence bool) {
	return r2.start-r1.end == 1
}

func (r *IPAddressV4Range) GetCountry() (c *country.Country) {
	if r.data == nil {
		return nil
	}
	return r.data.Country
}
