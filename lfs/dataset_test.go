package lfs

import (
	log "github.com/sirupsen/logrus"
	"math"
	conf "pds-go/lfs/config"
	"pds-go/lfs/dataset"
	"pds-go/lfs/db"
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
		t.Errorf("TestMean failed as mean value is incorrect, got: %f, want: %f", mean, 303.365663)
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
		t.Errorf("DropByColumn failed as NumColumns is incorrect, got: %d, want: %d.", cols, initial-1)
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
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}
