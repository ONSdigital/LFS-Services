package dataset_test

import (
	log "github.com/sirupsen/logrus"
	conf "lfs/lfs-services/config"
	"lfs/lfs-services/dataset"
	"lfs/lfs-services/db"
	"math"
	"testing"
)

func setupDataset(logger *log.Logger) (*dataset.Dataset, error) {
	d, err := dataset.NewDataset("test", logger)
	if err != nil {
		logger.Error(err)
		return &dataset.Dataset{}, nil
	}

	err = d.LoadSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", dataset.BigDataset{})
	if err != nil {
		logger.Error(err)
		return &dataset.Dataset{}, nil
	}

	return &d, nil
}

func TestMean(t *testing.T) {

	logger := log.New()

	ds, err := setupDataset(logger)
	if err != nil {
		panic(err)
	}

	mean, err := ds.Mean("QUOTA")
	if err != nil {
		panic(err)
	}

	v := math.Round(mean*100) / 100
	if v != 303.37 {
		t.Errorf("TestMean failed, got: %f, want: %f", mean, 303.365663)
	}
}

func TestDropColumn(t *testing.T) {
	logger := log.New()

	ds, err := setupDataset(logger)
	if err != nil {
		panic(err)
	}

	initial := ds.NumColumns()

	err = ds.DropColumn("QUOTA")
	if err != nil {
		panic(err)
	}

	cols := ds.NumColumns()
	if cols != initial-1 {
		t.Errorf("DropByColumn failed, got: %d, want: %d.", cols, initial-1)
	}
}

func TestFromSav(t *testing.T) {

	logger := log.New()

	d, err := dataset.NewDataset("test", logger)
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	err = d.LoadSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", dataset.BigDataset{})
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}
	logger.Printf("dataset contains %d row(s)\n", d.NumRows())
	_ = d.Head(5)
}

func TestToCSV(t *testing.T) {

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
	logger := log.New()

	d, err := dataset.NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	err = d.LoadSav(testDirectory()+"ips1710bv2.sav", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}

	err = d.ToCSV("out.csv")
	if err != nil {
		logger.Panic(err)
	}

	t.Logf("Dataset Size: %d\n", d.NumRows())
	_ = d.Head(5)
}

func TestToSav(t *testing.T) {
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

	logger := log.New()

	d, err := dataset.NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	err = d.LoadSav(testDirectory()+"ips1710bv2.sav", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}

	err = d.ToSAV(testDirectory() + "dataset-export.sav")
	if err != nil {
		logger.Panic(err)
	}
}

func TestFromCSV(t *testing.T) {

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

	logger := log.New()

	d, err := dataset.NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	err = d.LoadCSV(testDirectory()+"out.csv", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}

	logger.Printf("dataset contains %d row(s)\n", d.NumRows())
	_ = d.Head(5)
}

func TestUnPersist(t *testing.T) {
	logger := log.New()
	d, err := db.GetDefaultPersistenceImpl(logger).UnpersistDataset("LFSwk18PERS_non_confidential")
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}
	logger.Printf("dataset contains %d row(s)\n", d.NumRows())
	_ = d.Head(5)
}

func TestPersist(t *testing.T) {
	logger := log.New()

	d, err := dataset.NewDataset("test", logger)
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	err = d.LoadSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", dataset.BigDataset{})
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	err = db.GetDefaultPersistenceImpl(logger).PersistDataset(d)
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	logger.Printf("dataset contains %d row(s)\n", d.NumRows())
}

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in the configuration file")
	}
	return
}
