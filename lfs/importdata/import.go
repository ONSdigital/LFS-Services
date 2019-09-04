package importdata

import "pds-go/lfs/importdata/sav"

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
