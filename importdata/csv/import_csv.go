package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"strings"
)

type ImportCSV struct{}

func ImportCSVToSlice(fileName string) (out [][]string, err error) {
	csvfile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open CSV file: %s, error: %w", fileName, err)
	}
	defer func() { _ = csvfile.Close() }()

	r := csv.NewReader(csvfile)
	cnt := 0
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV file: %s, error: %w", fileName, err)
		}
		if cnt == 0 {
			for i, j := range record {
				record[i] = strings.ToUpper(j)
			}
		}
		cnt++
		out = append(out, record)
	}
	return out, nil
}

func (ImportCSV) Import(fileName string, out interface{}) error {
	clientsFile, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot open CSV file: %s, error: %s", fileName, err)
	}

	defer func() {
		_ = clientsFile.Close()
	}()

	if err := gocsv.UnmarshalFile(clientsFile, out); err != nil {
		return fmt.Errorf("cannot unmarshall CSV file: %s, error: %w", fileName, err)
	}

	return nil
}
