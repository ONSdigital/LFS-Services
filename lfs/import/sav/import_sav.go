package sav

// #cgo amd64 CFLAGS: -g
// #cgo LDFLAGS: -lreadstat
// #include <stdlib.h>
// #include "import_sav.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"log"
	"unsafe"
)

type Items struct {
	Shiftno float64 `json:"Shiftno"`
	Serial  float64 `json:"Serial"`
	Version string  `json:"Version"`
}

var items []Items

//export goAddLine
func goAddLine(str *C.char) {
	gostr := C.GoString(str)

	err := json.Unmarshal([]byte(gostr), &items)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Number of items: %d, size of buffer: %d, size struct %d\n", len(items), len(gostr), len(items))
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
