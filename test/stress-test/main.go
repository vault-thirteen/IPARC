package main

import (
	"math"
	"os"
	"time"

	"github.com/vault-thirteen/auxie/IPA"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/IPARC/common/helper"
	"github.com/vault-thirteen/IPARC/ipar"
	"github.com/vault-thirteen/IPARC/iparc"
)

func main() {
	if len(os.Args) != 2 {
		cio.LogMessage("db zip file is not set")
		os.Exit(1)
	}
	dbZipFile := os.Args[1]

	var err error
	var dbFilePath string
	dbFilePath, err = helper.UnpackDbFile(dbZipFile)
	mustBeNoError(err)
	defer func() {
		cio.LogMessage("deleting temporary folders ...")
		derr := helper.DeleteTemporaryDataFolders()
		if derr != nil {
			errors.LogError(derr)
		}
	}()

	var col *iparc.IPAddressV4RangeCollection
	col, err = iparc.NewFromCsvFile(dbFilePath)
	mustBeNoError(err)
	cio.LogMessage("Stress test has started. It may take several minutes. Please wait ...")

	t1 := time.Now()
	var addr ipa.IPAddressV4
	var rng *ipar.IPAddressV4Range
	for ba := 0; ba <= math.MaxUint8; ba++ {
		for bb := 0; bb <= math.MaxUint8; bb++ {
			for bc := 0; bc <= math.MaxUint8; bc++ {
				for bd := 0; bd <= math.MaxUint8; bd++ {
					addr = ipa.NewFromBytes(byte(ba), byte(bb), byte(bc), byte(bd))
					rng, err = col.GetRangeByIPAddress(addr)
					mustBeNoError(err)
				}
			}
		}
	}
	timeSpent := time.Now().Sub(t1).Seconds()
	var rps = float64(4294967295) / timeSpent
	cio.PrintMessage(cio.FormatMessage("Time spent: %v sec., RPS=%.2f.", timeSpent, rps))
	rng = rng
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
