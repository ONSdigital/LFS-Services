package exportdata_test

import (
	conf "pds-go/lfs/config"
	ext "pds-go/lfs/exportdata"
	imp "pds-go/lfs/importdata"
	"testing"
)

type TestDataset struct {
	Shiftno      float64 `spss:"Shiftno"`
	Serial       float64 `spss:"Serial"`
	Version      string  `spss:"Version"`
	PortRoute2   float64 `spss:"PortRoute2"`
	Baseport     string  `spss:"Baseport"`
	PRouteLatDeg float64 `spss:"PRouteLatDeg"`
	PRouteLonEW  string  `spss:"PRouteLonEW"`
	DVLineName   string  `spss:"DVLineName"`
	DVPortName   string  `spss:"DVPortName"`
}

func TestExportSav(t *testing.T) {

	var spssFile []TestDataset

	t.Logf("Reading sav")
	if err := imp.ImportSavFile(testDirectory()+"ips1710bv2.sav", &spssFile); err != nil {
		panic(err)
	}

	i := len(spssFile)

	t.Logf("Total Items: %d\n", i)

	t.Logf("Writing sav")
	if err := ext.ExportSavFile(testDirectory()+"export_test.sav", &spssFile); err != nil {
		t.Logf("Test failed - writer")
	} else {
		t.Logf("Test finished - writer")
	}
}

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}
