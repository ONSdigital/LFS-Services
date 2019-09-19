package sav

// #cgo windows amd64 CFLAGS: -g -IC:/msys64/mingw64/include
// #cgo windows LDFLAGS: -LC:/msys64/mingw64/lib -lreadstat
// #cgo darwin amd64 CFLAGS: -g
// #cgo darwin LDFLAGS: -lreadstat
// #cgo linux amd64 CFLAGS: -I/usr/local/include -g
// #cgo linux LDFLAGS: -L/usr/local/lib -lreadstat
// #include "sav_writer.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"pds-go/lfs/io/spss"
	"unsafe"
)

type Header struct {
	SavType spss.ColumnType
	Name    string
	Label   string
}

type DataItem struct {
	Value []interface{}
}

var DoubleValue []float64
var StringValue []string

func Export(fileName string, label string, headers []Header, data []DataItem) int {

	numHeaders := len(headers)

	cHeaders := (*[8192]*C.file_header)(C.malloc(C.size_t(C.sizeof_file_header * numHeaders)))

	for i, f := range headers {
		header := (*C.file_header)(C.malloc(C.size_t(C.sizeof_file_header)))
		(*header).sav_type = C.int(f.SavType)
		(*header).name = C.CString(f.Name)
		(*header).label = C.CString(f.Label)
		cHeaders[i] = header
	}

	numRows := len(data)
	// DataItem represents a single data item. The number of items is therefore the
	// number of rows multiplied by the number of columns
	cDataItem := (*[1 << 28]*C.data_item)(C.malloc(C.size_t(C.sizeof_data_item * numRows * numHeaders)))

	cnt := 0

	for _, r := range data {
		for j := 0; j < len(r.Value); j++ {
			col := r.Value[j]
			dataItem := (*C.data_item)(C.malloc(C.size_t(C.sizeof_data_item)))

			(*dataItem).sav_type = C.int(headers[j].SavType)

			switch headers[j].SavType {

			case spss.ReadstatTypeString:
				if _, ok := col.(string); !ok {
					(*dataItem).string_value = C.CString(col.(string))
					panic("Invalid type, string expected")
				}
				(*dataItem).string_value = C.CString(col.(string))

			case spss.ReadstatTypeInt8:
				if _, ok := col.(int); !ok {
					panic("Invalid type, int8 expected")
				}
				(*dataItem).int_value = C.int(col.(int))

			case spss.ReadstatTypeInt16:
				if _, ok := col.(int); !ok {
					panic("Invalid type, int16 expected")
				}
				(*dataItem).int_value = C.int(col.(int))

			case spss.ReadstatTypeInt32:
				if _, ok := col.(int); !ok {
					panic("Invalid type, int32 expected")
				}
				(*dataItem).int_value = C.int(col.(int))

			case spss.ReadstatTypeFloat:
				if _, ok := col.(float32); !ok {
					panic("Invalid type, float32 expected")
				}
				(*dataItem).float_value = C.float(col.(float32))

			case spss.ReadstatTypeDouble:
				if _, ok := col.(float64); !ok {
					fmt.Printf("Invalid type, double expected: %s\n", col)
					panic("Invalid type, double expected")
				}
				(*dataItem).double_value = C.double(col.(float64))

			case spss.ReadstatTypeStringRef:
				panic("String references not supported")
			}
			cDataItem[cnt] = dataItem
			cnt++
		}
	}

	res, err := C.save_sav(C.CString(fileName), C.CString(label), &cHeaders[0], C.int(numHeaders), C.int(numRows), &cDataItem[0])

	// Free up C allocated memory
	for i := 0; i < numHeaders; i++ {
		C.free(unsafe.Pointer((*cHeaders[i]).name))
		C.free(unsafe.Pointer((*cHeaders[i]).label))
		C.free(unsafe.Pointer(cHeaders[i]))
	}
	C.free(unsafe.Pointer(cHeaders))

	for i := 0; i < numRows*numHeaders; i++ {
		if int((*cDataItem[i]).sav_type) == spss.ReadstatTypeString.AsInt() {
			C.free(unsafe.Pointer((*cDataItem[i]).string_value))
		}
		C.free(unsafe.Pointer(cDataItem[i]))
	}
	C.free(unsafe.Pointer(cDataItem))

	if err != nil {
		fmt.Printf(" -> spss export: C code returned  %s", err)
	}

	return int(res)
}

func DefaultSPSSWriter(in interface{}) Writer {
	return FileOutput{in.(string)}
}

func WriteToSPSSFile(out string, in interface{}) error {
	return SpssWriter(out).Write(in)
}

var SpssWriter = DefaultSPSSWriter
