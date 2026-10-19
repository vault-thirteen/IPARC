package ipad

import (
	"github.com/vault-thirteen/IPARC/ipad/country"
)

type IPAddressData struct {
	*country.Country
}

func New(countryCode string, countryName string) (d *IPAddressData, err error) {
	d = new(IPAddressData)

	d.Country, err = country.New(countryCode, countryName)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *IPAddressData) CountryCode() string {
	return d.Country.Code()
}

func (d *IPAddressData) CountryName() string {
	return d.Country.Name()
}
