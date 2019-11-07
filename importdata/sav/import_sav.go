package sav

// #cgo windows amd64 CFLAGS: -O3 -IC:/msys64/mingw64/include
// #cgo windows LDFLAGS: -LC:/msys64/mingw64/libs -lreadstat
// #cgo darwin amd64 CFLAGS: -g
// #cgo darwin LDFLAGS: -lreadstat
// #cgo linux amd64 CFLAGS: -I/usr/local/include -g
// #cgo linux LDFLAGS: -L/usr/local/lib -lreadstat
// #include <stdlib.h>
// #include "sav_reader.h"
import "C"

import (
	"errors"
	"fmt"
	"os"
	"services/types"
	"unsafe"
)

func ImportSav(fileName string) (types.SavImportData, error) {

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return types.SavImportData{}, fmt.Errorf(" -> Import: file %s not found", fileName)
	}

	name := C.CString(fileName)
	defer C.free(unsafe.Pointer(name))

	var res = C.parse_sav(name)
	if res == nil {
		return types.SavImportData{}, errors.New("read from SPSS file failed")
	}

	defer func() {
		if res == nil {
			return
		}
		C.cleanup(res)
		C.free(unsafe.Pointer(res))
	}()

	v := C.struct_Data(*res)

	headerCount := int(v.header_count)
	rowCount := int(v.row_count)

	var header = make([]types.Header, headerCount)

	for i := 0; i < headerCount; i++ {
		var head **C.struct_Header = v.header
		z := (*[1 << 30]*C.struct_Header)((unsafe.Pointer(head)))[i]
		typeString := types.TypeString

		switch int(z.var_type) {
		case 0:
			typeString = types.TypeString
		case 1:
			typeString = types.TypeInt8
		case 2:
			typeString = types.TypeInt16
		case 3:
			typeString = types.TypeInt32
		case 4:
			typeString = types.TypeFloat
		case 5:
			typeString = types.TypeDouble
		}

		header[i] = types.Header{
			VariableName:        C.GoString(z.var_name),
			VariableDescription: C.GoString(z.var_description),
			VariableType:        typeString,
			VariableLength:      int(z.length),
			VariablePrecision:   int(z.precision),
		}
	}

	var savRows = make([]types.Rows, rowCount)

	for i := 0; i < rowCount; i++ {
		var rows **C.struct_Rows = v.rows
		r := (*[1 << 30]*C.struct_Rows)((unsafe.Pointer(rows)))[i]

		length := int(r.row_length)

		var rowValues = make([]string, length)

		rowData := r.row_data
		for j := 0; j < length; j++ {
			s := (*[1 << 30]*C.char)((unsafe.Pointer(rowData)))[j]
			rowValues[j] = C.GoString(s)
		}

		savRows[i] = types.Rows{RowData: rowValues}
	}

	savImportData := types.SavImportData{
		Header:      header,
		HeaderCount: headerCount,
		Rows:        savRows,
		RowCount:    rowCount,
	}

	return savImportData, nil
}

var spssReader = DefaultSPSSReader

type SavFileImport struct{}

func (SavFileImport) Import(fileName string, out interface{}) error {
	return spssReader(fileName).Read(out)
}

func DefaultSPSSReader(in interface{}) Reader {
	return FileInput{in.(string)}
}
