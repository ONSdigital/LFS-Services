package main

// #cgo amd64 CFLAGS: -g
// #cgo LDFLAGS: -lreadstat
// #include <stdlib.h>
// #include "sav_generate.h"
import "C"

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"strings"
	"unicode"
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

type HeaderInfo struct {
	Name    string
	VarType int
}

var header = make([]HeaderInfo, 0)

//export goAddHeaderItem
func goAddHeaderItem(pos C.int, name *C.char, varType C.int, end C.int) {
	header = append(header, HeaderInfo{strings.ToUpper(C.GoString(name)), int(varType)})
}

func Generate(fileName string) int {
	name := C.CString(fileName)
	defer C.free(unsafe.Pointer(name))

	res := C.read_header(name)
	if res != 0 {
		return 1
	}

	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	input := flag.String("input", "", "input file name")
	output := flag.String("output", "", "output file name")
	structName := flag.String("struct", "SpssDataItem", "structure name")
	packageName := flag.String("package", "lfs", "package name")

	flag.Parse()

	if *output == "" || *input == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Create(*output)
	check(err)

	defer f.Close()

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	Generate(*input)

	if len(header) == 0 {
		fmt.Println("No items found in SPSS file. Is it Empty?")
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(w, "package %s\n\n", *packageName)
	_, _ = fmt.Fprintf(w, "type %s struct {\n", *structName)

	for _, j := range header {
		a := []rune(j.Name)
		a[0] = unicode.ToUpper(a[0])
		name := string(a)
		switch j.VarType {
		case ReadstatTypeString:
			_, _ = fmt.Fprintf(w, "    %s string \t`spss:\"%s\"`\n", name, j.Name)
		case ReadstatTypeInt8:
			_, _ = fmt.Fprintf(w, "    %s int \t`spss:\"%s\"`\n", name, j.Name)
		case ReadstatTypeInt16:
			_, _ = fmt.Fprintf(w, "    %s int \t`spss:\"%s\"`\n", name, j.Name)
		case ReadstatTypeInt32:
			_, _ = fmt.Fprintf(w, "    %s int \t`spss:\"%s\"`\n", name, j.Name)
		case ReadstatTypeFloat:
			_, _ = fmt.Fprintf(w, "    %s float \t`spss:\"%s\"`\n", name, j.Name)
		case ReadstatTypeDouble:
			_, _ = fmt.Fprintf(w, "    %s float64 \t`spss:\"%s\"`\n", name, j.Name)
		case ReadstatTypeStringRef:
			panic("String references not supported")
		}
	}

	_, _ = fmt.Fprintf(w, "}\n")
	_ = w.Flush()

	content, err := format.Source(b.Bytes())

	_, _ = f.Write(content)
}
