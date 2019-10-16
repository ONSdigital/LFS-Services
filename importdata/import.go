package importdata

import (
	"services/importdata/csv"
	"services/importdata/sav"
)

type ImportFunction func(fileName string, out interface{}) error

type Importer interface {
	Import(fileName string, out interface{}) error
}

func importFile(i Importer) ImportFunction {
	return func(fileName string, out interface{}) error {
		return i.Import(fileName, out)
	}
}

var ImportSavFile = importFile(sav.SavFileImport{})
var ImportCSVFile = importFile(csv.ImportCSV{})
