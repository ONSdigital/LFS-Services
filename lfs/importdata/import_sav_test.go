package importdata_test

import (
	conf "pds-go/lfs/config"
	im "pds-go/lfs/importdata"
	"testing"
)

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}

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

func TestImportSav(t *testing.T) {

	var spssFile []TestDataset

	if err := im.ImportSavFile(testDirectory()+"ips1710bv2.sav", &spssFile); err != nil {
		panic(err)
	}

	t.Logf("Starting test - reader")

	i := len(spssFile)

	t.Logf("Total Items: %d\n", i)
	t.Logf("Test finished - reader")
}
