package sav

// #cgo amd64 CFLAGS: -g
// #cgo LDFLAGS: -lreadstat
// #include <stdlib.h>
// #include "import_sav.h"
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	ReadstatTypeString    = iota
	ReadstatTypeInt8      = iota
	ReadstatTypeInt16     = iota
	ReadstatTypeInt32     = iota
	ReadstatTypeFloat     = iota
	ReadstatTypeDouble    = iota
	ReadstatTypeStringRef = iota
)

type headerLine struct {
	name  string
	vType int
}

//type itemLine struct {
//	val string
//}

var headerItems = make(map[int]headerLine)

//var lineItems = make(map[int]itemLine)

//export goAddLine
func goAddLine(str *C.char) {
	gostr := C.GoString(str)
	println(gostr)
}

//export goAddHeaderLine
func goAddHeaderLine(pos C.int, name *C.char, varType C.int, end C.int) {
	if int(end) == 1 { // we are done
		printHeader()
	} else {
		headerItems[int(pos)] = headerLine{C.GoString(name), int(varType)}
	}
}

func printHeader() {
	for k := range headerItems {
		fmt.Printf("key[%d] -> title[%s] ", k, headerItems[k].name)
		switch headerItems[k].vType {
		case ReadstatTypeString:
			fmt.Printf(" -> type[string]\n")
		case ReadstatTypeInt8:
			fmt.Printf(" -> type[int8]\n")
		case ReadstatTypeInt16:
			fmt.Printf(" -> type[int16]\n")
		case ReadstatTypeInt32:
			fmt.Printf(" -> type[int32]\n")
		case ReadstatTypeFloat:
			fmt.Printf(" -> type[float]\n")
		case ReadstatTypeDouble:
			fmt.Printf(" -> type[double]\n")
		case ReadstatTypeStringRef:
			fmt.Printf(" -> type[string ref]\n")
		}
	}
}

func Import(fileName string) int {
	name := C.CString(fileName)
	defer C.free(unsafe.Pointer(name))

	res := C.parse_sav(name)
	if res != 0 {
		return 1
	}

	return 0
}
