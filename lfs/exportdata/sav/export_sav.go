package sav

type ExportSavFile struct{}

func (ExportSavFile) Export(out string, in interface{}) error {
	return SpssWriter(out).Write(in)
}
