package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/vault-thirteen/IPARC/common/helper"
	"github.com/vault-thirteen/IPARC/ipar"
	"github.com/vault-thirteen/IPARC/iparc"
	"github.com/vault-thirteen/auxie/IPA"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("db zip file is not set")
	}
	dbZipFile := os.Args[1]

	var err error
	var dbFilePath string
	dbFilePath, err = helper.UnpackDbFile(dbZipFile)
	mustBeNoError(err)
	defer func() {
		log.Println("deleting temporary folders ...")
		derr := helper.DeleteTemporaryDataFolders()
		if derr != nil {
			log.Println(derr)
		}
	}()

	var col *iparc.IPAddressV4RangeCollection
	col, err = iparc.NewFromCsvFile(dbFilePath)
	mustBeNoError(err)
	log.Println("Stress test has started. It may take several minutes. Please wait ...")

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
	fmt.Printf("Time spent: %v sec., RPS=%.2f.\r\n", timeSpent, rps)
	rng = rng
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
