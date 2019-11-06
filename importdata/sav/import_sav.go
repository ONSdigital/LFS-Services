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

const EOL = "\n"

func ImportSav(fileName string) ([][]string, error) {

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, fmt.Errorf(" -> Import: file %s not found", fileName)
	}

	name := C.CString(fileName)
	defer C.free(unsafe.Pointer(name))

	var res = C.parse_sav(name)
	if res == nil {
		return nil, errors.New("read from SPSS file failed")
	}

	var str [][]string

	defer func() {
		if res == nil {
			return
		}
		C.cleanup(res)
		C.free(unsafe.Pointer(res))
	}()

	v := C.struct_Data(*res)

	j := int(v.header_count)

	var header = make([]types.Definitions, j)
	for i := 0; i < j; i++ {
		var head **C.struct_Header = v.header
		z := (*[1 << 30]*C.struct_Header)((unsafe.Pointer(head)))[i]
		type_string := types.TYPE_STRING

		switch int(z.var_type) {
		case 0:
			type_string = types.TYPE_STRING
		case 1:
			type_string = types.TYPE_INT8
		case 2:
			type_string = types.TYPE_INT16
		case 3:
			type_string = types.TYPE_INT32
		case 4:
			type_string = types.TYPE_FLOAT
		case 5:
			type_string = types.TYPE_DOUBLE
		}
		header[i] = types.Definitions{
			Variable:       C.GoString(z.var_name),
			Description:    C.GoString(z.var_description),
			VariableType:   type_string,
			VariableLength: int(z.length),
			Precision:      int(z.precision),
			Alias:          "",
			Editable:       false,
			Imputation:     false,
			DV:             false,
		}
	}

	for _, j := range header {
		fmt.Printf("Var: %s, \tDesc: %s, \tType: %s, \tLength: %d, \tPrecision: %d, \tAlias: %s\n", j.Variable, j.Description, j.VariableType,
			j.VariableLength, j.Precision, j.Alias)
	}
	//header := []string{C.GoString(v.header)}

	//for _, l := range header {
	//	s := strings.Split(l, typ.TagSeparator)
	//	str = append(str, s)
	//}

	//data := strings.Split(C.GoString(v.data), EOL)
	//
	//for _, l := range data {
	//	s := strings.Split(l, typ.TagSeparator)
	//	str = append(str, s)
	//}

	return str, nil
}

var spssReader = DefaultSPSSReader

type SavFileImport struct{}

func (SavFileImport) Import(fileName string, out interface{}) error {
	return spssReader(fileName).Read(out)
}

func DefaultSPSSReader(in interface{}) Reader {
	return FileInput{in.(string)}
}
