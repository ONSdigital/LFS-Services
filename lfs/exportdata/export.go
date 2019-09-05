package exportdata

import "pds-go/lfs/exportdata/sav"
import "pds-go/lfs/exportdata/csv"

type ExportFunction func(fileName string, out interface{}) error

type ExportData interface {
	Export(out string, in interface{}) error
}

func exportFile(i ExportData) ExportFunction {
	return func(fileName string, out interface{}) error {
		return i.Export(fileName, out)
	}
}

var ExportSavFile = exportFile(sav.ExportSavFile{})
var ExportCSVFile = exportFile(csv.ExportCSVFile{})
