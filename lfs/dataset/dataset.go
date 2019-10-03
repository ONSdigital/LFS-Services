package dataset

import "C"
import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"pds-go/lfs/importdata/sav"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Column struct {
	ColNo int
	Kind  reflect.Kind
	Rows  []interface{}
}

type Dataset struct {
	TableName   string
	Columns     map[string]Column
	mux         *sync.Mutex
	logger      *log.Logger
	RowCount    int
	ColumnCount int
}

const (
	InitialRowCapacity    = 20000
	InitialColumnCapacity = 2000
	COMMA                 = ","
)

func NewDataset(name string, logger *log.Logger) (Dataset, error) {
	mux := sync.Mutex{}
	cols := make(map[string]Column, InitialColumnCapacity)
	return Dataset{name, cols, &mux, logger, 0, 0}, nil
}

type fromFileFunc func(fileName string, out interface{}) error

func (d *Dataset) logTime(from fromFileFunc) fromFileFunc {
	return func(fileName string, out interface{}) error {
		startTime := time.Now()
		err := from(fileName, out)
		a := time.Now().Sub(startTime)
		d.logger.Printf("load processed in %s", a)
		return err
	}
}

func (d *Dataset) LoadSav(fileName string, out interface{}) error {
	return d.logTime(d.readSav)(fileName, out)
}

func (d *Dataset) readSav(in string, out interface{}) error {
	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		return fmt.Errorf(" -> FromSav: %T is not a struct type", out)
	}

	start := time.Now()

	records, err := sav.ImportSav(in)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf(" -> populateDataset: spss file: %s is empty", in)
	}

	elapsed := time.Since(start)
	d.logger.Info(fmt.Sprintf("read sav file (%d records) in %s", len(records)-1, elapsed))

	start = time.Now()
	er := d.populateDataset(in, records, out)
	if er != nil {
		return er
	}
	elapsed = time.Since(start)
	d.logger.Info(fmt.Sprintf("created dataset (%d records) in %s", d.RowCount, elapsed))

	return nil
}

func (d *Dataset) AddRow(row map[string]interface{}) error {
	//d.mux.Lock()
	//defer d.mux.Unlock()

	if len(row) != len(d.Columns) {
		return fmt.Errorf(" -> AddRow: Column count mismatch. Expected %d, got %d", len(d.Columns), len(row))
	}
	for k, v := range row {
		col := d.Columns[k]
		col.Rows = append(col.Rows, v)
		d.Columns[k] = col
	}
	d.RowCount++
	return nil
}

func (d *Dataset) AddColumn(name string, columnType reflect.Kind) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.Columns[name]; ok {
		return fmt.Errorf(" -> AddColumn: Column %s already exists", name)
	}

	col := Column{}
	col.Kind = columnType
	col.ColNo = d.ColumnCount
	col.Rows = make([]interface{}, 0, InitialRowCapacity)
	d.Columns[name] = col
	d.ColumnCount++

	if d.RowCount == 0 {
		return nil
	}

	// Add empty Rows if we have existing data
	for i := 0; i < d.RowCount; i++ {
		switch columnType {
		case reflect.String:
			col.Rows = append(col.Rows, "")
		case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
			col.Rows = append(col.Rows, 0)
		case reflect.Float32, reflect.Float64:
			col.Rows = append(col.Rows, 0.0)
		default:
			return fmt.Errorf(" -> AddColumn: cannot convert type")
		}
	}

	return nil
}

func (d *Dataset) DropColumn(name string) error {
	if _, ok := d.Columns[name]; !ok {
		return fmt.Errorf(" -> DropColumn: Column %s does not exist", name)
	}

	a := d.OrderedColumns()
	m := make(map[string]Column, InitialColumnCapacity)

	var colNo = 0
	for _, v := range a {
		if v != name {
			var col Column
			col.Rows = m[v].Rows
			col.Kind = m[v].Kind
			col.ColNo = colNo
			m[v] = col
			colNo++
		}
	}

	d.Columns = m
	d.ColumnCount--
	return nil
}

func (d *Dataset) populateDataset(fileName string, rows [][]string, out interface{}) error {
	_, file := filepath.Split(fileName)
	var extension = filepath.Ext(file)
	var name = file[0 : len(file)-len(extension)]
	var err error
	*d, err = NewDataset(name, d.logger)

	if err != nil {
		return fmt.Errorf(" -> populateDataset: cannot create a new DataSet: %s", err)
	}

	d.logger.Info("starting import into Dataset")

	t1 := reflect.TypeOf(out)

	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		if err := d.AddColumn(strings.ToUpper(a.Name), a.Type.Kind()); err != nil {
			return fmt.Errorf(" -> populateDataset: cannot create column: %w", err)
		}
	}

	headers := rows[0]
	body := rows[1:]

	for a := range headers {
		headers[a] = strings.ToUpper(headers[a])
	}

	for _, spssRow := range body {

		row := make(map[string]interface{})

		for j := 0; j < len(spssRow); j++ {
			if len(spssRow) != len(headers) {
				return fmt.Errorf(" -> populateDataset: header is out of alignment with row. row size: %d, column size: %d", len(spssRow), len(headers))
			}
			header := strings.ToUpper(headers[j])
			// extract the tagged columns only
			if _, ok := d.Columns[headers[j]]; !ok {
				continue
			}

			// check type is valid
			a := spssRow[j]
			if a == "" {
				a = "NULL"
			}

			kind := d.Columns[headers[j]].Kind
			switch kind {

			case reflect.String:
				row[header] = a

			case reflect.Int8, reflect.Uint8:
				i, err := strconv.ParseInt(a, 0, 8)
				if err != nil {
					return fmt.Errorf(" -> populateDataset: cannot convert %s into an Int8", a)
				}
				row[header] = i

			case reflect.Int, reflect.Int32, reflect.Uint32:
				i, err := strconv.ParseInt(a, 0, 32)
				if err != nil {
					return fmt.Errorf(" -> populateDataset: cannot convert %s into an Int32", a)
				}
				row[header] = i

			case reflect.Int64, reflect.Uint64:
				i, err := strconv.ParseInt(a, 0, 64)
				if err != nil {
					return fmt.Errorf(" -> populateDataset: cannot convert %s into an Int64", a)
				}
				row[header] = i

			case reflect.Float32:
				i, err := strconv.ParseFloat(a, 32)
				if err != nil {
					return fmt.Errorf(" -> populateDataset: cannot convert %s into an Float32", a)
				}
				row[header] = i

			case reflect.Float64:
				i, err := strconv.ParseFloat(a, 64)
				if err != nil {
					return fmt.Errorf(" -> populateDataset: cannot convert %s into an Float64", a)
				}
				row[header] = i

			default:
				return fmt.Errorf(" -> populateDataset: cannot convert struct variable type from SPSS type")
			}
		}

		if err := d.AddRow(row); err != nil {
			return fmt.Errorf(" -> populateDataset: AddRow failed %w", err)
		}

	}
	return nil
}

func (d Dataset) OrderedColumns() []string {
	var keys = make([]string, d.ColumnCount)
	for k, v := range d.Columns {
		keys[v.ColNo] = k
	}
	return keys
}

func (d *Dataset) getAllRows() ([]string, [][]string) {
	return d.getByRow(d.RowCount, d.ColumnCount)
}

func (d *Dataset) getByRow(maxRows int, maxCols int) ([]string, [][]string) {
	cnt := 0
	var header []string
	var items [][]string

	if maxCols > d.ColumnCount {
		maxCols = d.ColumnCount
	}

	for _, v := range d.OrderedColumns() {
		if cnt > maxCols-1 {
			break
		}
		header = append(header, v)
		cnt++
	}

	if maxRows > d.RowCount {
		maxRows = d.RowCount
	}

	// for each header, get MaxRows
	for j := 0; j < maxRows; j++ {
		var row []string
		for _, b := range header {
			r := d.Columns[b].Rows[j]
			kind := d.Columns[b].Kind

			switch kind {
			case reflect.String:
				row = append(row, r.(string))
			case reflect.Int8, reflect.Uint8:
				row = append(row, fmt.Sprintf("%d", r.(int)))
			case reflect.Int, reflect.Int32, reflect.Uint32:
				row = append(row, fmt.Sprintf("%d", r.(int64)))
			case reflect.Int64, reflect.Uint64:
				row = append(row, fmt.Sprintf("%d", r.(int64)))
			case reflect.Float32:
				row = append(row, fmt.Sprintf("%f", r.(float32)))
			case reflect.Float64:
				row = append(row, fmt.Sprintf("%g", r.(float64)))
			default:
				panic(fmt.Errorf(" -> getByRow: unknown type - possible corruption"))
			}
		}
		items = append(items, row)
	}
	return header, items
}

func (d *Dataset) Head(max ...int) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	var maxItems = 5
	if max != nil {
		maxItems = max[0]
	}

	table := tablewriter.NewWriter(os.Stdout)

	const maxCols = 15
	header, items := d.getByRow(maxItems, maxCols)

	table.SetHeader(header)
	for _, b := range items {
		table.Append(b)
	}

	j := fmt.Sprintf("%d Rows(s)\n", table.NumLines())
	table.SetCaption(true, j)
	table.Render()
	return nil
}

func (d Dataset) NumColumns() int {
	return d.ColumnCount
}

func (d Dataset) NumRows() int {
	return d.RowCount
}

func (d Dataset) Mean(name string) (float64, error) {
	if _, ok := d.Columns[name]; !ok {
		return 0.0, fmt.Errorf(" -> Mean: Column %s does not exist", name)
	}

	var kind = d.Columns[name].Kind

	if kind == reflect.String {
		return 0.0, errors.New(fmt.Sprintf(" -> Mean: column %s is not numeric", name))
	}

	var avg = 0.0

	for _, v := range d.Columns[name].Rows {
		switch kind {
		case reflect.Int8, reflect.Uint8:
			avg = avg + float64(v.(int))
		case reflect.Int, reflect.Int32, reflect.Uint32:
			avg = avg + float64(v.(int))
		case reflect.Int64, reflect.Uint64:
			avg = avg + float64(v.(int))
		case reflect.Float32:
			avg = avg + float64(v.(float32))
		case reflect.Float64:
			avg = avg + v.(float64)
		default:
			panic(fmt.Errorf(" -> Mean: unknown type - possible corruption"))
		}
	}

	return avg / float64(d.RowCount), nil
}
