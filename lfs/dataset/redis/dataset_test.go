package redis

import (
	log "github.com/sirupsen/logrus"
	"math"
	conf "pds-go/lfs/config"
	"testing"
)

func setupDataset(logger *log.Logger) (*Dataset, error) {
	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Error(err)
		return &Dataset{}, nil
	}

	err = d.LoadSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", BigDataset{})
	if err != nil {
		logger.Error(err)
		return &Dataset{}, nil
	}

	return &d, nil
}

func TestMean(t *testing.T) {

	logger := log.New()

	dataset, err := setupDataset(logger)
	if err != nil {
		panic(err)
	}

	mean, err := dataset.Mean("QUOTA")
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

	dataset, err := setupDataset(logger)
	if err != nil {
		panic(err)
	}

	initial := dataset.NumColumns()

	err = dataset.DropColumn("QUOTA")
	if err != nil {
		panic(err)
	}

	cols := dataset.NumColumns()
	if cols != initial-1 {
		t.Errorf("DropByColumn failed as NumColumns is incorrect, got: %d, want: %d.", cols, initial-1)
	}
}

func TestFromSav(t *testing.T) {

	logger := log.New()

	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	err = d.LoadSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", BigDataset{})
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}
	logger.Printf("dataset contains %d row(s)\n", d.NumRows())
	_ = d.Head(5)
}

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}
