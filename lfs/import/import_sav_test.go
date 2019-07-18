package importdata

import (
	"fmt"
)

import "testing"

func TestImportSav(t *testing.T) {
	res := ImportSavFile("testdata/ips1710bv2.sav")

	if res != 0 {
		t.Error("Received error from parse_sav")
		t.FailNow()
	}

	fmt.Println("Done")
}
