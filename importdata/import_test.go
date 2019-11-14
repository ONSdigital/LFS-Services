package importdata_test

import (
	conf "services/config"
	im "services/importdata"
	"services/importdata/sav"
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

	res, err := sav.ImportSav(testDirectory() + "LFSwkJANUARYNI_non_confidential.SAV")
	if err != nil {
		panic(err)
	}

	headerCount := res.HeaderCount
	rowCount := res.RowCount
	labelsCount := res.LabelsCount

	t.Logf("Total Columns: %d\n", headerCount)
	t.Logf("Total Rows: %d\n", rowCount)
	t.Logf("Total Value Labels: %d\n", labelsCount)
	t.Logf("Test finished - spss reader")
}
