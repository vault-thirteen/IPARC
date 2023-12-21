package iparc

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/vault-thirteen/IPARC/ipad"
	"github.com/vault-thirteen/IPARC/ipad/country"
	"github.com/vault-thirteen/IPARC/ipar"
	"github.com/vault-thirteen/auxie/IPA"
	"github.com/vault-thirteen/auxie/as"
	ae "github.com/vault-thirteen/auxie/errors"
)

const (
	ErrSyntaxErrorOnLine  = "syntax error on line %d: %s"
	ErrNotEnoughColumns   = "not enough columns"
	ErrStrangeCountryCode = "strange country code: %s"
	ErrStrangeCountryName = "strange country name: %s"
	ErrIPAddress          = "IP address error: %s"
	ErrCountry            = "country error: %s"
	ErrSequence           = "sequence error, range: %+v"
	ErrData               = "data error: %v"

	// ErrAdjacentRangeIsNotFound error means inconsistency of the database.
	// It may only happen when not all IP addresses are listed in the
	// collection / database.
	ErrAdjacentRangeIsNotFound = "adjacent range is not found, idx=%v, ipa=%v"
)

type IPAddressV4RangeCollection struct {
	ranges    []*ipar.IPAddressV4Range
	countries map[string]*country.Country

	// index is an array for fast search.
	// Each item is a middle point of a range.
	// Indices in 'index' array are the same as in 'ranges' array.
	index []float64
}

func NewFromCsvFile(filePath string) (col *IPAddressV4RangeCollection, err error) {
	col = new(IPAddressV4RangeCollection)
	err = col.init()
	if err != nil {
		return nil, err
	}

	err = col.fillWithDataFromCsvFile(filePath)
	if err != nil {
		return nil, err
	}

	col.createIndex()

	return col, nil
}

func (col *IPAddressV4RangeCollection) init() (err error) {
	col.ranges = make([]*ipar.IPAddressV4Range, 0)
	col.countries = make(map[string]*country.Country)

	var unknownCountry *country.Country
	unknownCountry, err = country.New(country.CodeUnknown, country.NameUnknown)
	if err != nil {
		return err
	}

	col.countries[country.CodeUnknown] = unknownCountry

	return nil
}

func (col *IPAddressV4RangeCollection) fillWithDataFromCsvFile(filePath string) (err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = true
	var csvRecord []string
	var lineNumber = 1
	var rng *ipar.IPAddressV4Range
	for {
		csvRecord, err = csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		rng, err = col.processCsvLine(csvRecord)
		if err != nil {
			return fmt.Errorf(ErrSyntaxErrorOnLine, lineNumber, err.Error())
		}

		err = col.addIPARange(rng)
		if err != nil {
			return err
		}

		lineNumber++
	}

	return nil
}

func (col *IPAddressV4RangeCollection) processCsvLine(csvRecord []string) (rng *ipar.IPAddressV4Range, err error) {
	if len(csvRecord) < 4 {
		return nil, errors.New(ErrNotEnoughColumns)
	}

	var ipaStart ipa.IPAddressV4
	ipaStart, err = ipa.NewFromUintString(csvRecord[0])
	if err != nil {
		return nil, fmt.Errorf(ErrIPAddress, err.Error())
	}

	var ipaEnd ipa.IPAddressV4
	ipaEnd, err = ipa.NewFromUintString(csvRecord[1])
	if err != nil {
		return nil, fmt.Errorf(ErrIPAddress, err.Error())
	}

	var countryCode, countryName string
	countryCode = csvRecord[2]
	countryName = csvRecord[3]
	if countryCode == country.CodeUnknown {
		countryName = country.NameUnknown
	} else {
		if !country.MayBeCountryCode(countryCode) {
			return nil, fmt.Errorf(ErrStrangeCountryCode, countryCode)
		}
		if len(countryName) == 0 {
			return nil, fmt.Errorf(ErrStrangeCountryName, countryName)
		}
	}

	var cnt *country.Country
	cnt, err = col.findOrCreateCountry(countryCode, countryName)
	if err != nil {
		return nil, fmt.Errorf(ErrCountry, err.Error())
	}

	var data *ipad.IPAddressData
	data, err = ipad.New(cnt.Code(), cnt.Name())
	if err != nil {
		return nil, fmt.Errorf(ErrData, err.Error())
	}

	return ipar.New(ipaStart, ipaEnd, data)
}

func (col *IPAddressV4RangeCollection) findOrCreateCountry(
	countryCode string,
	countryName string,
) (c *country.Country, err error) {
	if countryName == country.NameUnknown {
		return col.countries[country.CodeUnknown], nil
	}

	var countryExists bool
	c, countryExists = col.countries[countryCode]
	if countryExists {
		return c, nil
	}

	c, err = country.New(countryCode, countryName)
	if err != nil {
		return nil, err
	}
	col.countries[countryCode] = c
	return c, nil
}

func (col *IPAddressV4RangeCollection) addIPARange(rng *ipar.IPAddressV4Range) (err error) {
	if len(col.ranges) == 0 {
		col.ranges = append(col.ranges, rng)
		return nil
	}

	lastIPARange := col.ranges[len(col.ranges)-1]
	if !ipar.IsSequence(lastIPARange, rng) {
		return fmt.Errorf(ErrSequence, rng)
	}

	col.ranges = append(col.ranges, rng)
	return nil
}

func (col *IPAddressV4RangeCollection) createIndex() {
	col.index = make([]float64, 0, len(col.ranges))
	for _, rng := range col.ranges {
		col.index = append(col.index, rng.GetMiddle())
	}
}

func (col *IPAddressV4RangeCollection) GetRangeByIPAddress(
	ipa ipa.IPAddressV4,
) (r *ipar.IPAddressV4Range, err error) {
	v := float64(ipa)
	midIndexA := col.searchNearestMiddle(v)
	rngA := col.ranges[midIndexA]
	midValueA := rngA.GetMiddle()
	radiusA := rngA.GetRadius()

	// Check the radius.
	var adjacentRangeIndex int
	if midValueA < v {
		if v <= midValueA+radiusA {
			return rngA, nil
		} else {
			// Adjacent right range is our result.
			// If the database is full, it should exist.
			adjacentRangeIndex = midIndexA + 1
			if adjacentRangeIndex > len(col.ranges)-1 {
				err = fmt.Errorf(ErrAdjacentRangeIsNotFound, adjacentRangeIndex, ipa)
				return nil, err
			}
			return col.ranges[adjacentRangeIndex], nil
		}
	} else if v < midValueA {
		if midValueA-radiusA <= v {
			return rngA, nil
		} else {
			// Adjacent left range is our result.
			// If the database is full, it should exist.
			adjacentRangeIndex = midIndexA - 1
			if adjacentRangeIndex < 0 {
				err = fmt.Errorf(ErrAdjacentRangeIsNotFound, adjacentRangeIndex, ipa)
				return nil, err
			}
			return col.ranges[adjacentRangeIndex], nil
		}
	} else { // v == midValueA.
		return rngA, nil
	}
}

func (col *IPAddressV4RangeCollection) searchNearestMiddle(v float64) (midIndex int) {
	indices := as.FindNearestBS(col.index, v)

	// For sequential ranges with a gap of 1.0 the result may have two indices
	// returned only if v was not a round number (not integer). This is not
	// possible because IPv4 address is always a round number (integer). So, we
	// can safely ignore cases of two middle points returned by the search.

	return indices[0]
}
