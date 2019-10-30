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
	Width   int
}

var header = make([]HeaderInfo, 0)

//export goAddHeaderItem
func goAddHeaderItem(pos C.int, name *C.char, varType C.int, end C.int, width C.int) {
	if int(end) != 1 { // we are done
		header = append(header, HeaderInfo{C.GoString(name), int(varType), int(width)})
	}
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

func generateStruct(input, output, packageName, structName *string) {
	f, err := os.Create(*output)
	check(err)

	defer func() { _ = f.Close() }()

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
			_, _ = fmt.Fprintf(w, "    %s string \t`spss:\"%s\" db:\"%s\"`\n", name, j.Name, name)
		case ReadstatTypeInt8:
			_, _ = fmt.Fprintf(w, "    %s int \t`spss:\"%s\" db:\"%s\"`\n", name, j.Name, name)
		case ReadstatTypeInt16:
			_, _ = fmt.Fprintf(w, "    %s int \t`spss:\"%s\" db:\"%s\"`\n", name, j.Name, name)
		case ReadstatTypeInt32:
			_, _ = fmt.Fprintf(w, "    %s int \t`spss:\"%s\" db:\"%s\"`\n", name, j.Name, name)
		case ReadstatTypeFloat:
			_, _ = fmt.Fprintf(w, "    %s float \t`spss:\"%s\" db:\"%s\"`\n", name, j.Name, name)
		case ReadstatTypeDouble:
			_, _ = fmt.Fprintf(w, "    %s float64 \t`spss:\"%s\" db:\"%s\"`\n", name, j.Name, name)
		case ReadstatTypeStringRef:
			panic("String references not supported")
		}
	}

	_, _ = fmt.Fprintf(w, "}\n")
	_ = w.Flush()

	content, err := format.Source(b.Bytes())

	_, _ = f.Write(content)
}

func generateTable(input *string, output *string, tableName *string) {
	f, err := os.Create(*output)
	check(err)

	defer func() { _ = f.Close() }()

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	Generate(*input)

	if len(header) == 0 {
		fmt.Println("No items found in SPSS file. Is it Empty?")
		os.Exit(1)
	}

	_, _ = fmt.Fprintf(w, "create table %s\n", *tableName)
	_, _ = fmt.Fprintf(w, "(\n")
	_, _ = fmt.Fprintf(w, "    id int not null,\n")
	_, _ = fmt.Fprintf(w, "    file_name varchar(64) not null,\n")
	_, _ = fmt.Fprintf(w, "    file_source char(2) not null,\n")
	_, _ = fmt.Fprintf(w, "    week int not null,\n")
	_, _ = fmt.Fprintf(w, "    month int not null,\n")
	_, _ = fmt.Fprintf(w, "    year int not null,\n")

	for i, j := range header {
		a := []rune(j.Name)
		a[0] = unicode.ToUpper(a[0])
		name := string(a)
		switch j.VarType {
		case ReadstatTypeString:
			_, _ = fmt.Fprintf(w, "    %s varchar(%d) sparse null", name, j.Width)
		case ReadstatTypeInt8:
			_, _ = fmt.Fprintf(w, "    %s tinyint sparse null", name)
		case ReadstatTypeInt16:
			_, _ = fmt.Fprintf(w, "    %s smallint sparse null", name)
		case ReadstatTypeInt32:
			_, _ = fmt.Fprintf(w, "    %s int sparse null", name)
		case ReadstatTypeFloat:
			_, _ = fmt.Fprintf(w, "    %s float sparse null", name)
		case ReadstatTypeDouble:
			_, _ = fmt.Fprintf(w, "    %s real sparse null", name)
		case ReadstatTypeStringRef:
			panic("String references not supported")
		}
		if i != len(header)-1 {
			_, _ = fmt.Fprintf(w, ",\n")
		} else {
			_, _ = fmt.Fprintf(w, "\n")
		}
	}

	_, _ = fmt.Fprintf(w, ")\n")
	_ = w.Flush()

	_, err = f.Write(b.Bytes())
}

func main() {

	input := flag.String("input", "", "input file name")
	output := flag.String("output", "", "output file name")

	structName := flag.String("struct", "SpssDataItem", "structure name")
	packageName := flag.String("package", "lfs", "package name")
	spssOutput := flag.Bool("gen-struct", false, "generate Go struct")

	tableOutput := flag.Bool("gen-table", false, "generate create table DDL")
	tableName := flag.String("tableName", "table", "table name in DDL")

	flag.Parse()

	if *output == "" || *input == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *spssOutput && *tableOutput {
		fmt.Println("either -gen-struct or -gen-table but not both")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !*spssOutput && !*tableOutput {
		fmt.Println("select one of either -gen-struct or -gen-table")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *spssOutput {
		generateStruct(input, output, packageName, structName)
	} else {
		generateTable(input, output, tableName)
	}

}
