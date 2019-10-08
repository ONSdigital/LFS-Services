package exportdata_test

import (
	conf "services/config"
	ext "services/exportdata"
	imp "services/importdata"
	"testing"
)

func TestExportCSV(t *testing.T) {
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

	var csv []TestDataset

	t.Logf("Reading csv")
	if err := imp.ImportCSVFile(testDirectory()+"out.csv", &csv); err != nil {
		panic(err)
	}

	i := len(csv)

	t.Logf("Total Items: %d\n", i)

	t.Logf("Writing csv")

	if err := ext.ExportCSVFile(testDirectory()+"export_test.csv", &csv); err != nil {
		t.Logf("Test failed - csv writer")
	} else {
		t.Logf("Test finished - csv writer")
	}
}

func TestExportSav(t *testing.T) {

	type TestDataset struct {
		Shiftno      float64 `spss:"Shiftno"`
		Serial       string  `spss:"Serial"`
		Version      string  `spss:"Version"`
		PortRoute2   float64 `spss:"PortRoute2"`
		Baseport     string  `spss:"Baseport"`
		PRouteLatDeg float64 `spss:"PRouteLatDeg"`
		PRouteLonEW  string  `spss:"PRouteLonEW"`
		DVLineName   string  `spss:"DVLineName"`
		DVPortName   string  `spss:"DVPortName"`
	}

	var spssFile []TestDataset

	t.Logf("Reading sav")
	if err := imp.ImportSavFile(testDirectory()+"ips1710bv2.sav", &spssFile); err != nil {
		panic(err)
	}

	i := len(spssFile)

	t.Logf("Total Items: %d\n", i)

	t.Logf("Writing sav")
	if err := ext.ExportSavFile(testDirectory()+"export_test.sav", &spssFile); err != nil {
		t.Logf("Test failed - sav writer")
	} else {
		t.Logf("Test finished - sav writer")
	}
}

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}
