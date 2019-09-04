package exportdata

import "pds-go/lfs/exportdata/sav"

type savFileExport struct{}

func (savFileExport) Export(out string, in interface{}) error {
	return sav.SpssWriter(out).Write(in)
}

type ExportFunction func(fileName string, out interface{}) error

type ExportData interface {
	Export(out string, in interface{}) error
}

func exportFile(i ExportData) ExportFunction {
	return func(fileName string, out interface{}) error {
		return i.Export(fileName, out)
	}
}

var ExportSavFile = exportFile(savFileExport{})
