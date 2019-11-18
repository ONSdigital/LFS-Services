package sav

// #include <stdlib.h>
// #include "sav_reader.h"
// #cgo windows amd64 CFLAGS: -IC:/msys64/mingw64/include
// #cgo windows LDFLAGS: -LC:/msys64/mingw64/lib -lreadstat
// #cgo darwin amd64 CFLAGS: -g
// #cgo darwin LDFLAGS: -lreadstat
// #cgo linux amd64 CFLAGS: -I/usr/local/include -g
// #cgo linux LDFLAGS: -L/usr/local/lib -lreadstat
import "C"

import (
	"errors"
	"fmt"
	"os"
	"services/types"
	"strings"
	"unsafe"
)

func ImportSav(fileName string) (types.SavImportData, error) {

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return types.SavImportData{}, fmt.Errorf(" -> Import: file %s not found", fileName)
	}

	name := C.CString(fileName)

	var res = C.parse_sav(name)
	if res == nil {
		return types.SavImportData{}, errors.New("read from sav file failed")
	}

	defer func() {
		C.free(unsafe.Pointer(name))
		if res == nil {
			return
		}
		C.cleanup(res)
		C.free(unsafe.Pointer(res))
	}()

	v := C.struct_Data(*res)

	headerCount := int(v.header_count)
	labelsCount := int(v.labels_count)
	rowCount := int(v.row_count)

	// Get headers
	var header = make([]types.Header, headerCount)
	var head **C.struct_Header = v.header

	for i := 0; i < headerCount; i++ {
		z := (*[1 << 30]*C.struct_Header)((unsafe.Pointer(head)))[i]
		typeString := getType(int(z.var_type))

		// Convert all variable names to uppercase
		header[i] = types.Header{
			VariableName:        strings.ToUpper(C.GoString(z.var_name)),
			VariableDescription: C.GoString(z.var_description),
			VariableType:        typeString,
			VariableLength:      int(z.length),
			VariablePrecision:   int(z.precision),
			LabelName:           strings.ToUpper(C.GoString(z.label_name)),
			Drop:                false,
		}
	}

	// get labels
	labelsMap := make(map[string][]types.Labels, labelsCount)
	var labels **C.struct_Labels = v.labels
	for i := 0; i < labelsCount; i++ {
		z := (*[1 << 30]*C.struct_Labels)((unsafe.Pointer(labels)))[i]
		name := strings.ToUpper(C.GoString(z.name))
		typeString := getType(int(z.var_type))
		var value interface{} = 0

		switch typeString {
		case types.TypeString:
			value = C.GoString(z.string_value)
		case types.TypeInt8:
			value = int64(z.i8_value)
		case types.TypeInt16:
			value = int64(z.i16_value)
		case types.TypeInt32:
			value = int64(z.i32_value)
		case types.TypeFloat:
			value = float32(z.float_value)
		case types.TypeDouble:
			value = float64(z.double_value)
		}

		label := C.GoString(z.label)

		ll := types.Labels{
			Name:          strings.ToUpper(name),
			Value:         value,
			Label:         label,
			Tag:           rune(z.tag),
			TagMissing:    int(z.tag_missing),
			SystemMissing: int(z.system_missing),
			VariableType:  typeString,
		}

		labelSet, ok := labelsMap[name]
		if !ok {
			valueA := make([]types.Labels, 0)
			valueA = append(valueA, ll)
			labelsMap[name] = valueA
		} else {
			labelSet = append(labelSet, ll)
			labelsMap[name] = labelSet
		}
	}

	// get rows
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
		Labels:      labelsMap,
		LabelsCount: labelsCount,
		Rows:        savRows,
		RowCount:    rowCount,
	}

	return savImportData, nil
}

func getType(savType int) types.SavType {
	typeString := types.TypeString

	switch savType {
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

	return typeString
}

var spssReader = DefaultSPSSReader

type SPSSFileImport struct{}

func (SPSSFileImport) Import(fileName string, out interface{}) error {
	return spssReader(fileName).Read(out)
}

func DefaultSPSSReader(in interface{}) Reader {
	return FileInput{in.(string)}
}
