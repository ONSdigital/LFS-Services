package csv

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type ExportCSVFile struct{}

func (ExportCSVFile) Export(fileName string, out interface{}) error {

	clientsFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = clientsFile.Close()
	}()

	if err := gocsv.MarshalFile(out, clientsFile); err != nil {
		return fmt.Errorf("cannot marshall CSV file: %s, err: %w", clientsFile, err)
	}

	return nil
}
