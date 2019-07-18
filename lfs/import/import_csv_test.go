package importdata

import (
	"fmt"
	"testing"
)

func TestImportCsv(t *testing.T) {
	res := ImportCsvFile("testdata/ips1710bv2.sav")

	if res != 0 {
		t.Error("Received error from parse_csv")
	}

	fmt.Println("Done")
}
