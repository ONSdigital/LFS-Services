package importdata_test

import (
	conf "services/config"
	im "services/importdata"
	"testing"
)

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}

func TestImportCSV(t *testing.T) {

	type TestDataset struct {
		Shiftno      float64 `csv:"Shiftno"`
		Serial       float64 `csv:"Serial"`
		Version      string  `csv:"Version"`
		PortRoute2   float64 `csv:"PortRoute2"`
		Baseport     string  `csv:"Baseport"`
		PRouteLatDeg float64 `csv:"PRouteLatDeg"`
		PRouteLonEW  string  `csv:"PRouteLonEW"`
		DVLineName   string  `csv:"DVLineName"`
		DVPortName   string  `csv:"DVPortName"`
	}

	var csvFile []TestDataset

	if err := im.ImportCSVFile(testDirectory()+"out.csv", &csvFile); err != nil {
		panic(err)
	}

	t.Logf("Starting test - csv reader")

	i := len(csvFile)

	t.Logf("Total Items: %d\n", i)
	t.Logf("Test finished - csv reader")
}

func TestImportSav(t *testing.T) {

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

	var spssFile []TestDataset

	if err := im.ImportSavFile(testDirectory()+"ips1710bv2.sav", &spssFile); err != nil {
		panic(err)
	}

	t.Logf("Starting test - spss reader")

	i := len(spssFile)

	t.Logf("Total Items: %d\n", i)
	t.Logf("Test finished - spss reader")
}
