package importdata

import (
	"lfs/lfs-services/importdata/csv"
	"lfs/lfs-services/importdata/sav"
)

type ImportFunction func(fileName string, out interface{}) error

type ImportData interface {
	Import(fileName string, out interface{}) error
}

func importFile(i ImportData) ImportFunction {
	return func(fileName string, out interface{}) error {
		return i.Import(fileName, out)
	}
}

var ImportSavFile = importFile(sav.SavFileImport{})
var ImportCSVFile = importFile(csv.ImportCSV{})
