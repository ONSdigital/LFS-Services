package exportdata

import "services/exportdata/sav"
import "services/exportdata/csv"

type ExportFunction func(fileName string, out interface{}) error

type Exporter interface {
	Export(out string, in interface{}) error
}

func exportFile(i Exporter) ExportFunction {
	return func(fileName string, out interface{}) error {
		return i.Export(fileName, out)
	}
}

var ExportSavFile = exportFile(sav.ExportSavFile{})
var ExportCSVFile = exportFile(csv.ExportCSVFile{})
