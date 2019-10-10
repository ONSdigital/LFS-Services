package dataset

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	di "services/exportdata/sav"
	imcsv "services/importdata/csv"
	"services/importdata/sav"
	"services/io/spss"
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
	DatasetName string
	Columns     map[string]Column
	mux         *sync.Mutex
	RowCount    int
	ColumnCount int
}

const (
	InitialRowCapacity    = 20000
	InitialColumnCapacity = 2000
)

func NewDataset(name string) (Dataset, error) {
	mux := sync.Mutex{}
	cols := make(map[string]Column, InitialColumnCapacity)
	return Dataset{name, cols, &mux, 0, 0}, nil
}

type fromFileFunc func(fileName, datasetName string, out interface{}) error

func (d *Dataset) logTime(from fromFileFunc) fromFileFunc {
	return func(fileName, datasetName string, out interface{}) error {
		startTime := time.Now()
		err := from(fileName, datasetName, out)
		a := time.Now().Sub(startTime)

		log.WithFields(log.Fields{
			"method":      "logTime",
			"file":        fileName,
			"elapsedTime": a,
		}).Debug("Load processed")

		return err
	}
}

func (d *Dataset) LoadCSV(fileName, datasetName string, out interface{}) error {
	return d.logTime(d.readCSV)(fileName, datasetName, out)
}

func (d *Dataset) readCSV(in, datasetName string, out interface{}) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		log.WithFields(log.Fields{
			"method": "readCSV",
			"file":   in,
		}).Error("The output interface is not a struct")
		return fmt.Errorf(" -> FromCSV: %T is not a struct type", out)
	}

	start := time.Now()
	records, err := imcsv.ImportCSVToSlice(in)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "readCSV",
			"file":   in,
		}).Error("Cannot import CSV file")
		return fmt.Errorf(" -> FromCSV: cannot import CSV file %w", err)
	}

	if len(records) == 0 {
		log.WithFields(log.Fields{
			"method": "readCSV",
			"file":   in,
		}).Warn("The CSV file is empty")
		return fmt.Errorf(" -> FromCSV: csv file: %s is empty", in)
	}

	elapsed := time.Since(start)
	log.WithFields(log.Fields{
		"method":      "readCSV",
		"file":        in,
		"records":     len(records) - 1,
		"elapsedTime": elapsed,
	}).Debug("Read CSV file")

	start = time.Now()
	err = d.populateDataset(in, datasetName, records, out)
	if err != nil {
		return err
	}

	elapsed = time.Since(start)

	log.WithFields(log.Fields{
		"method":      "readCSV",
		"records":     d.RowCount,
		"elapsedTime": elapsed,
	}).Debug("Dataset created")

	return nil
}

func (d *Dataset) LoadSav(fileName, datasetName string, out interface{}) error {
	return d.logTime(d.readSav)(fileName, datasetName, out)
}

func (d *Dataset) readSav(in, datasetName string, out interface{}) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		log.WithFields(log.Fields{
			"method": "readSav",
			"file":   in,
		}).Error("The output interface is not a struct")
		return fmt.Errorf(" -> readSav: %T is not a struct type", out)
	}

	start := time.Now()

	records, err := sav.ImportSav(in)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		log.WithFields(log.Fields{
			"method": "readSav",
			"file":   in,
		}).Warn("The sav file is empty")
		return fmt.Errorf(" -> readSav: spss file: %s is empty", in)
	}

	elapsed := time.Since(start)

	log.WithFields(log.Fields{
		"method":      "readSav",
		"file":        in,
		"records":     len(records) - 1,
		"elapsedTime": elapsed,
	}).Debug("Read sav file")

	start = time.Now()
	er := d.populateDataset(in, datasetName, records, out)
	if er != nil {
		return er
	}
	elapsed = time.Since(start)

	log.WithFields(log.Fields{
		"method":      "readSav",
		"records":     d.RowCount,
		"elapsedTime": elapsed,
	}).Debug("Dataset created")

	return nil
}

func (d Dataset) GetRows(colName string) ([]interface{}, error) {
	col, ok := d.Columns[colName]
	if !ok {
		return nil, fmt.Errorf("column %s not found", colName)
	}

	return col.Rows, nil
}

// if we had generics, this would not have to be repeated for each type....
func (d Dataset) GetRowsAsString(colName string) ([]string, error) {
	r, err := d.GetRows(colName)
	if err != nil {
		return nil, err
	}

	if d.Columns[colName].Kind != reflect.String {
		return nil, fmt.Errorf("column %s is not a string", colName)
	}

	rows := make([]string, d.RowCount)
	for _, a := range r {
		rows = append(rows, a.(string))
	}
	return rows, nil
}

func (d Dataset) GetRowsAsInt(colName string) ([]int, error) {

	r, err := d.GetRows(colName)
	if err != nil {
		return nil, err
	}

	if d.Columns[colName].Kind != reflect.Int {
		return nil, fmt.Errorf("column %s is not an int", colName)
	}

	rows := make([]int, d.RowCount)
	for _, a := range r {
		rows = append(rows, a.(int))
	}
	return rows, nil
}

func (d Dataset) GetRowsAsFloat(colName string) ([]float32, error) {

	r, err := d.GetRows(colName)
	if err != nil {
		return nil, err
	}

	if d.Columns[colName].Kind != reflect.Float32 {
		return nil, fmt.Errorf("column %s is not a float32", colName)
	}

	rows := make([]float32, d.RowCount)
	for _, a := range r {
		rows = append(rows, a.(float32))
	}
	return rows, nil
}

func (d Dataset) GetRowsAsDouble(colName string) ([]float64, error) {

	r, err := d.GetRows(colName)
	if err != nil {
		return nil, err
	}

	if d.Columns[colName].Kind != reflect.Float64 {
		return nil, fmt.Errorf("column %s is not a float64t", colName)
	}

	rows := make([]float64, d.RowCount)
	for _, a := range r {
		rows = append(rows, a.(float64))
	}
	return rows, nil
}

func (d Dataset) ToSAV(fileName string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	var header []di.Header
	var cols = d.OrderedColumns()

	for i := 0; i < len(cols); i++ {
		var spssType spss.ColumnType = 0

		switch d.Columns[cols[i]].Kind {
		case reflect.String:
			spssType = spss.ReadstatTypeString
		case reflect.Int8, reflect.Uint8:
			spssType = spss.ReadstatTypeInt8
		case reflect.Int, reflect.Int32, reflect.Uint32:
			spssType = spss.ReadstatTypeInt32
		case reflect.Float32:
			spssType = spss.ReadstatTypeFloat
		case reflect.Float64:
			spssType = spss.ReadstatTypeDouble
		default:
			log.WithFields(log.Fields{
				"method":   "ToSAV",
				"variable": cols[i],
			}).Error("Cannot convert type for struct variable into equivelent SPSS type")
			return fmt.Errorf("cannot convert type for struct variable %s into equivelent SPSS type", cols[i])
		}
		header = append(header, di.Header{SavType: spssType, Name: cols[i], Label: cols[i] + " label"})
	}

	h, items := d.getAllRows()
	var data []di.DataItem

	for _, v := range items {
		var dataItem di.DataItem
		dataItem.Value = make([]interface{}, 0)

		for j, value := range v {
			kind := d.Columns[h[j]].Kind
			switch kind {
			case reflect.String:
				dataItem.Value = append(dataItem.Value, fmt.Sprintf("%s", value))
			case reflect.Int8, reflect.Uint8:
				cv, err := strconv.ParseInt(value, 0, 8)
				if err != nil {
					return fmt.Errorf(" -> toSAV: cannot convert %s into an Int8", value)
				}
				dataItem.Value = append(dataItem.Value, cv)
			case reflect.Int, reflect.Int32, reflect.Uint32:
				cv, err := strconv.ParseInt(value, 0, 32)
				if err != nil {
					return fmt.Errorf(" -> toSAV: cannot convert %s into an Int32", value)
				}
				dataItem.Value = append(dataItem.Value, cv)
			case reflect.Int64, reflect.Uint64:
				cv, err := strconv.ParseInt(value, 0, 64)
				if err != nil {
					return fmt.Errorf(" -> toSAV: cannot convert %s into an Int64", value)
				}
				dataItem.Value = append(dataItem.Value, cv)
			case reflect.Float32:
				cv, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return fmt.Errorf(" -> toSAV: cannot convert %s into a Float32", value)
				}
				dataItem.Value = append(dataItem.Value, cv)
			case reflect.Float64:
				cv, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf(" -> toSAV: cannot convert %s into a Float64", value)
				}
				dataItem.Value = append(dataItem.Value, cv)
			default:
				return fmt.Errorf(" -> ToSAV: unknown type - possible corruption")
			}
		}
		data = append(data, dataItem)
	}

	if val := di.Export(fileName, d.DatasetName, header, data); val != 0 {
		log.WithFields(log.Fields{"method": "ToSAV",
			"file": fileName,
		}).Error("SPSS export failed")
		return fmt.Errorf(" -> ToSAV: spss export to %s failed", fileName)
	}

	return nil
}

func (d Dataset) ToCSV(fileName string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	f, err := os.Create(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "ToCSV",
			"file":   fileName,
		}).Error("Cannot create CSV output file")
		return fmt.Errorf(" -> ToCSV: cannot open output csv file: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	header, items := d.getAllRows()
	var buffer bytes.Buffer

	for i := 0; i < len(header); i++ {
		j := fmt.Sprintf("%s", header[i])
		buffer.WriteString(j)
		if i != len(header)-1 {
			buffer.WriteString(",")
		} else {
			buffer.WriteString("\n")
		}
	}

	q := buffer.String()

	if _, err = f.WriteString(q); err != nil {
		log.WithFields(log.Fields{
			"method": "ToCSV",
			"file":   fileName,
		}).Error("Cannot write to CSV file")
		return fmt.Errorf(" -> ToCSV: write to file: %s failed: %s", fileName, err)
	}

	for _, v := range items {
		buffer.Reset()
		for j, value := range v {
			buffer.WriteString(fmt.Sprintf("%s", value))
			if j != len(header)-1 {
				buffer.WriteString(",")
			} else {
				buffer.WriteString("\n")
			}
		}

		q := buffer.String()

		if _, err = f.WriteString(q); err != nil {
			log.WithFields(log.Fields{
				"method": "ToCSV",
				"file":   fileName,
			}).Error("Cannot write to CSV file")
			return fmt.Errorf(" -> ToCSV: write to file: %s failed: %s", fileName, err)
		}
	}

	return nil
}

func (d *Dataset) AddRow(row map[string]interface{}) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if len(row) != len(d.Columns) {
		log.WithFields(log.Fields{
			"method":   "AddRow",
			"expected": len(d.Columns),
			"got":      len(row),
		}).Error("Column count mismatch")
		return fmt.Errorf("column count mismatch. Expected %d, got %d", len(d.Columns), len(row))
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
		log.WithFields(log.Fields{
			"method": "AddColumn",
			"column": name,
		}).Warn("Column already exists")
		return fmt.Errorf("column %s already exists", name)
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
			log.WithFields(log.Fields{
				"method":     "AddColumn",
				"columnName": name,
				"columnType": columnType,
			}).Error("Cannot convert type")
			return fmt.Errorf("cannot convert type")
		}
	}

	return nil
}

func (d *Dataset) RenameColumns(columns map[string]string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	a := d.OrderedColumns()
	m := make(map[string]Column, InitialColumnCapacity)

	var colNo = 0
	for _, v := range a {
		colName := v
		if to, ok := columns[colName]; ok {
			colName = to
			log.WithFields(log.Fields{
				"from": v,
				"to":   colName,
			}).Debug("Rename column")
		}

		var col Column
		old := d.Columns[v]
		col.Rows = old.Rows
		col.Kind = old.Kind
		col.ColNo = colNo
		m[colName] = col
		colNo++
	}

	d.Columns = m
	return nil
}

func (d *Dataset) RenameColumn(from, to string) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	if _, ok := d.Columns[from]; !ok {
		log.WithFields(log.Fields{
			"method":     "RenameColumn",
			"fromColumn": from,
			"toColumn":   to,
		}).Warn("Column doesn't exist")
		return fmt.Errorf("column %s does not exist", from)
	}

	a := d.OrderedColumns()
	m := make(map[string]Column, InitialColumnCapacity)

	var colNo = 0
	for _, v := range a {
		colName := v
		if v == from {
			colName = to
		}

		var col Column
		old := d.Columns[v]
		col.Rows = old.Rows
		col.Kind = old.Kind
		col.ColNo = colNo
		m[colName] = col
		colNo++
	}

	d.Columns = m
	return nil
}

func isInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (d *Dataset) DropColumns(columns []string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	a := d.OrderedColumns()
	m := make(map[string]Column, InitialColumnCapacity)

	var colNo = 0
	for _, v := range a {
		if !isInSlice(v, columns) {
			var col Column
			old := d.Columns[v]
			col.Rows = old.Rows
			col.Kind = old.Kind
			col.ColNo = colNo
			m[v] = col
			colNo++
		} else {
			log.WithFields(log.Fields{
				"columnName": v,
			}).Debug("Dropping column")
		}
	}

	d.Columns = m
	d.ColumnCount = colNo
	return nil
}

func (d *Dataset) DropColumn(name string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.Columns[name]; !ok {
		log.WithFields(log.Fields{
			"method":     "DropColumn",
			"columnName": name,
		}).Warn("Column doesn't exist")
		return fmt.Errorf("column %s does not exist", name)
	}

	a := d.OrderedColumns()
	m := make(map[string]Column, InitialColumnCapacity)

	var colNo = 0
	for _, v := range a {
		if v != name {
			var col Column
			old := d.Columns[v]
			col.Rows = old.Rows
			col.Kind = old.Kind
			col.ColNo = colNo
			m[v] = col
			colNo++
		}
	}

	d.Columns = m
	d.ColumnCount--
	return nil
}

func (d *Dataset) populateDataset(fileName, datasetName string, rows [][]string, out interface{}) error {

	var err error
	*d, err = NewDataset(datasetName)

	if err != nil {
		log.WithFields(log.Fields{
			"datasetName":  datasetName,
			"errorMessage": err.Error(),
		}).Error("Cannot create dataset")
		return fmt.Errorf("cannot create a new DataSet: %s", err)
	}

	log.WithFields(log.Fields{
		"datasetName": datasetName,
	}).Debug("Starting import into Dataset")

	t1 := reflect.TypeOf(out)

	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		if err := d.AddColumn(strings.ToUpper(a.Name), a.Type.Kind()); err != nil {
			log.WithFields(log.Fields{
				"datasetName":  datasetName,
				"methodName":   "populateDataset",
				"columnName":   strings.ToUpper(a.Name),
				"columnType":   a.Type.Kind(),
				"errorMessage": err.Error(),
			}).Error("Cannot create column")
			return fmt.Errorf("cannot create column: %w", err)
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
				log.WithFields(log.Fields{
					"methodName": "populateDataset",
					"rowSize":    len(spssRow),
					"columnSize": len(headers),
				}).Error("header is out of alignment with row")
				return fmt.Errorf("header is out of alignment with row")
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
					logStructError("populateDataset", header, kind, "Int8")
					return fmt.Errorf("cannot convert %s into an Int8", a)
				}
				row[header] = i

			case reflect.Int, reflect.Int32, reflect.Uint32:
				i, err := strconv.ParseInt(a, 0, 32)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int32")
					return fmt.Errorf("cannot convert %s into an Int32", a)
				}
				row[header] = i

			case reflect.Int64, reflect.Uint64:
				i, err := strconv.ParseInt(a, 0, 64)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int64")
					return fmt.Errorf("cannot convert %s into an Int64", a)
				}
				row[header] = i

			case reflect.Float32:
				i, err := strconv.ParseFloat(a, 32)
				if err != nil {
					logStructError("populateDataset", header, kind, "Float32")
					return fmt.Errorf("cannot convert %s into an Float32", a)
				}
				row[header] = i

			case reflect.Float64:
				i, err := strconv.ParseFloat(a, 64)
				if err != nil {
					logStructError("populateDataset", header, kind, "Float64")
					return fmt.Errorf("cannot convert %s into an Float64", a)
				}

				row[header] = i

			default:
				logStructError("populateDataset", header, kind, "Unknown")
				return fmt.Errorf("cannot convert struct variable type from SPSS type")
			}
		}

		if err := d.AddRow(row); err != nil {
			log.WithFields(log.Fields{"methodName": "populateDataset"}).Error("Cannot add row")
			return fmt.Errorf("cannot add a row: %w", err)
		}

	}
	return nil
}

func logStructError(methodName, variableName string, kind reflect.Kind, newType string) {
	log.WithFields(log.Fields{
		"methodName":  methodName,
		"variable":    variableName,
		"convertFrom": kind,
		"convertTo":   newType,
	}).Error("Camnnot convert type")
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
				log.WithFields(log.Fields{"methodName": "getByRow", "type": kind}).Error("Unknown type - possible corruption")
				panic(fmt.Errorf("unknown type - possible corruption"))
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
		log.WithFields(log.Fields{"methodName": "Mean", "columnName": name}).Warn("Column does not exist")
		return 0.0, fmt.Errorf("column %s does not exist", name)
	}

	var kind = d.Columns[name].Kind

	if kind == reflect.String {
		log.WithFields(log.Fields{"methodName": "Mean", "columnName": kind}).Warn("column is not numeric")
		return 0.0, errors.New(fmt.Sprintf("column %s is not numeric", name))
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
			log.WithFields(log.Fields{"methodName": "Mean", "type": kind}).Error("Unknown type - possible corruption")
			return 0.0, fmt.Errorf("unknown type - possible corruption")
		}
	}

	return avg / float64(d.RowCount), nil
}
