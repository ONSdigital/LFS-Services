package redis

import (
	log "github.com/sirupsen/logrus"
	conf "pds-go/lfs/config"
	"testing"
)

func TestFromSav(t *testing.T) {

	logger := log.New()

	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

	_, err = d.FromSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", BigDataset{})
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}

}

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}
