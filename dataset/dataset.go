package dataset

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"math"
	"os"
	"reflect"
	di "services/exportdata/sav"
	imcsv "services/importdata/csv"
	"services/importdata/sav"
	"services/io/spss"
	"services/types"
	"services/util"
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
	types.Audit
}

const (
	InitialRowCapacity    = 20000
	InitialColumnCapacity = 2000
)

func NewDataset(name string) (Dataset, error) {
	mux := sync.Mutex{}
	cols := make(map[string]Column, InitialColumnCapacity)
	return Dataset{name, cols, &mux, 0, 0, types.Audit{}}, nil
}

func (d *Dataset) LoadCSV(fileName, datasetName string, out interface{}, filter types.Filter) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		log.Error().
			Str("method", "readCSV").
			Msg("The output interface is not a struct")
		return fmt.Errorf(" -> FromCSV: %T is not a struct type", out)
	}

	start := time.Now()
	records, err := imcsv.ImportCSVToSlice(fileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "readCSV").
			Str("file", fileName).
			Msg("Cannot import CSV file")
		return fmt.Errorf("cannot import CSV file %w", err)
	}

	if len(records) < 2 {
		log.Warn().
			Str("method", "readCSV").
			Msg("The CSV file is empty")
		return fmt.Errorf("csv file: %s is empty", fileName)
	}

	log.Info().
		Str("method", "readCSV").
		Str("file", fileName).
		Int("records", len(records)-1).
		Str("elapsedTime", util.FmtDuration(start)).
		Msg("Read CSV file")

	start = time.Now()
	err = d.populateDataset(datasetName, records, out, filter)
	if err != nil {
		return err
	}

	log.Debug().
		Str("method", "readCSV").
		Str("file", "in").
		Int("records", d.NumRows()).
		Str("elapsedTime", util.FmtDuration(start)).
		Msg("Populated Dataset")

	return nil
}

func (d *Dataset) LoadSav(in string, datasetName string, out interface{}, filter types.Filter) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		log.Error().
			Str("method", "readSav").
			Msg("The output interface is not a struct")
		return fmt.Errorf("%T is not a struct type", out)
	}

	start := time.Now()

	records, err := sav.ImportSav(in)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		log.Warn().
			Str("method", "readSav").
			Msg("The CSV file is empty")
		return fmt.Errorf("the spss file: %s is empty", in)
	}

	log.Debug().
		Str("file", in).
		Int("records", len(records)-1).
		Str("elapsedTime", util.FmtDuration(start)).
		Msg("Read Sav file")

	start = time.Now()
	er := d.populateDataset(datasetName, records, out, filter)
	if er != nil {
		return er
	}

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

	rows := make([]string, 0, d.NumRows())
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

	rows := make([]int, 0, d.NumRows())
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

	rows := make([]float32, 0, d.NumRows())
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

	rows := make([]float64, 0, d.NumRows())
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
			log.Error().
				Str("method", "ToSAV").
				Str("variable", cols[i]).
				Msg("Cannot convert type for struct variable into equivelent SPSS type")
			return fmt.Errorf("cannot convert type for struct variable %s into equivelent SPSS type", cols[i])
		}
		header = append(header, di.Header{SavType: spssType, Name: cols[i], Label: cols[i] + " label"})
	}

	h, items := d.GetAllRows()
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
		log.Error().
			Str("file", fileName).
			Msg("SPSS export failed")
		return fmt.Errorf(" -> ToSAV: spss export to %s failed", fileName)
	}

	return nil
}

func (d Dataset) ToCSV(fileName string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	f, err := os.Create(fileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "ToCSV").
			Str("file", fileName).
			Msg("Cannot create CSV output file")
		return fmt.Errorf(" -> ToCSV: cannot open output csv file: %s", err)
	}

	defer func() {
		_ = f.Close()
	}()

	header, items := d.GetAllRows()
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
		log.Error().
			Err(err).
			Str("method", "ToCSV").
			Str("file", fileName).
			Msg("Cannot write to CSV file")
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
			log.Error().
				Err(err).
				Str("file", fileName).
				Msg("Cannot write to CSV file")
			return fmt.Errorf(" -> ToCSV: write to file: %s failed: %s", fileName, err)
		}
	}

	return nil
}

func (d *Dataset) AddRow(row map[string]interface{}) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if len(row) != len(d.Columns) {
		log.Error().
			Str("method", "AddRow").
			Str("expected", string(len(d.Columns))).
			Str("got", string(len(row))).
			Msg("Column count mismatch")

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

func (d *Dataset) AddColumn(name string, columnType reflect.Kind) (*Column, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.Columns[name]; ok {
		log.Warn().
			Str("method", "AddColumn").
			Str("column", name).
			Msg("Column already exists")
		return &Column{}, fmt.Errorf("column %s already exists", name)
	}

	col := Column{}
	col.Kind = columnType
	col.ColNo = d.ColumnCount
	col.Rows = make([]interface{}, 0, InitialRowCapacity)
	d.ColumnCount++

	if d.RowCount == 0 {
		d.Columns[name] = col
		return &col, nil
	}

	// Add empty Rows if we have existing data
	for i := 0; i < d.RowCount; i++ {
		switch columnType {
		case reflect.String:
			col.Rows = append(col.Rows, "")
		case reflect.Int8:
			col.Rows = append(col.Rows, int8(0))
		case reflect.Uint8:
			col.Rows = append(col.Rows, uint8(0))
		case reflect.Int:
			col.Rows = append(col.Rows, int(0))
		case reflect.Int32:
			col.Rows = append(col.Rows, int32(0))
		case reflect.Uint32:
			col.Rows = append(col.Rows, uint32(0))
		case reflect.Int64:
			col.Rows = append(col.Rows, int64(0))
		case reflect.Uint64:
			col.Rows = append(col.Rows, uint64(0))
		case reflect.Float32:
			col.Rows = append(col.Rows, float32(0.0))
		case reflect.Float64:
			col.Rows = append(col.Rows, float64(0.0))
		default:
			log.Error().
				Str("method", "AddColumn").
				Str("column", name).
				Str("columnType", string(columnType)).
				Msg("Cannot convert type")
			return &Column{}, fmt.Errorf("cannot convert type")
		}
	}

	d.Columns[name] = col

	return &col, nil
}

func (d *Dataset) RenameColumn(from, to string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.Columns[from]; !ok {
		log.Warn().
			Str("method", "RenameColumns").
			Str("fromColumn", from).
			Str("toColumn", to).
			Msg("Column doesn't exist")
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

func (d *Dataset) DropColumn(name string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.Columns[name]; !ok {
		log.Warn().
			Str("method", "DropColumn").
			Str("columnName", name).
			Msg("Column doesn't exist")

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

func (d *Dataset) populateDataset(datasetName string, rows [][]string, out interface{}, filter types.Filter) error {

	var err error
	*d, err = NewDataset(datasetName)

	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot create a dataset")
		return fmt.Errorf("cannot create a  DataSet: %s", err)
	}

	log.Debug().
		Str("datasetName", datasetName).
		Msg("Starting import into Dataset")

	t1 := reflect.TypeOf(out)

	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		// skip columns that are marked as being dropped
		if filter.DropColumn(strings.ToUpper(a.Name)) {
			continue
		}
		if _, err := d.AddColumn(strings.ToUpper(a.Name), a.Type.Kind()); err != nil {
			log.Error().
				Err(err).
				Str("datasetName", datasetName).
				Str("methodName", "populateDataset").
				Str("columnName", strings.ToUpper(a.Name)).
				Str("columnType", string(a.Type.Kind())).
				Msg("Cannot create column")
			return fmt.Errorf("cannot create column: %w", err)
		}
	}

	headers := rows[0]
	body := rows[1:]

	d.ReferenceDate = time.Now()
	d.NumObFile = len(body)
	d.NumObLoaded = len(body)
	d.NumVarFile = len(headers)
	d.NumVarLoaded = len(headers)

	for a := range headers {
		headers[a] = strings.ToUpper(headers[a])
	}

	for _, spssRow := range body {

		row := make(map[string]interface{})

		for j := 0; j < len(spssRow); j++ {
			if len(spssRow) != len(headers) {
				log.Error().
					Err(err).
					Str("methodName", datasetName).
					Str("methodName", "populateDataset").
					Str("rowSize", string(len(spssRow))).
					Str("columnSize", string(len(headers))).
					Msg("Header is out of alignment with row")
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
				if d.Columns[headers[j]].Kind == reflect.String {
					a = "NULL"
				} else {
					a = "NaN"
				}
			}

			kind := d.Columns[headers[j]].Kind

			switch kind {

			case reflect.String:
				row[header] = a

			case reflect.Int8, reflect.Uint8:
				if a == "NaN" {
					a = "0"
				}
				i, err := strconv.ParseInt(a, 0, 8)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int8")
					return fmt.Errorf("cannot convert %s into an Int8", a)
				}
				row[header] = i

			case reflect.Int, reflect.Int32, reflect.Uint32:
				if a == "NaN" {
					a = "0"
				}
				i, err := strconv.ParseInt(a, 0, 32)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int32")
					return fmt.Errorf("cannot convert %s into an Int32", a)
				}
				row[header] = i

			case reflect.Int64, reflect.Uint64:
				if a == "NaN" {
					a = "0"
				}
				i, err := strconv.ParseInt(a, 0, 64)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int64")
					return fmt.Errorf("cannot convert %s into an Int64", a)
				}
				row[header] = i

			case reflect.Float32:
				if a == "NaN" {
					row[header] = math.NaN()
				} else {
					i, err := strconv.ParseFloat(a, 32)

					if err != nil {
						logStructError("populateDataset", header, kind, "Float32")
						return fmt.Errorf("cannot convert %s into an Float32", a)
					}
					row[header] = i
				}

			case reflect.Float64:
				if a == "NaN" {
					row[header] = math.NaN()
				} else {
					i, err := strconv.ParseFloat(a, 64)
					if err != nil {
						logStructError("populateDataset", header, kind, "Float64")
						return fmt.Errorf("cannot convert %s into an Float64", a)
					}

					row[header] = i
				}

			default:
				logStructError("populateDataset", header, kind, "Unknown")
				return fmt.Errorf("cannot convert struct variable type from SPSS type")
			}
		}

		// call the skipRow filter
		if filter.SkipRow(row) {
			continue
		}

		if err := d.AddRow(row); err != nil {
			log.Error().
				Err(err).
				Str("methodName", "populateDataset").
				Msg("Camnnot add row")
			return fmt.Errorf("cannot add a row: %w", err)
		}
	}

	m := make(map[string]Column, d.NumColumns())

	for k, v := range d.Columns {
		to, ok := filter.RenameColumns(k)
		if ok {
			m[to] = v
		} else {
			m[k] = v
		}
	}

	d.Columns = m

	return nil
}

func logStructError(methodName, variableName string, kind reflect.Kind, newType string) {
	log.Error().
		Str("methodName", methodName).
		Str("variable", variableName).
		Str("convertFrom", string(kind)).
		Str("convertTo", string(newType)).
		Msg("Camnnot convert type")
}

func (d Dataset) OrderedColumns() []string {
	var keys = make([]string, d.ColumnCount)
	for k, v := range d.Columns {
		keys[v.ColNo] = k
	}
	return keys
}

func (d *Dataset) GetAllRows() (headers []string, rows [][]string) {
	headers, rows = d.getByRow(d.RowCount, d.ColumnCount)
	return
}

func (d *Dataset) getByRow(maxRows int, maxCols int) ([]string, [][]string) {

	if maxCols > d.NumColumns() {
		maxCols = d.NumColumns()
	}

	var header = make([]string, maxCols)

	for k, v := range d.OrderedColumns() {
		if k > maxCols-1 {
			break
		}
		header[k] = v
	}

	if maxRows > d.NumRows() {
		maxRows = d.NumRows()
	}
	var items = make([][]string, maxRows)

	// for each header, get MaxRows
	for j := 0; j < maxRows; j++ {
		var row = make([]string, 0, 2000)
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
				log.Error().
					Str("methodName", "getByRow").
					Msg("Unknown type - possible corruption")
				panic(fmt.Errorf("unknown type - possible corruption"))
			}
		}
		items[j] = row
		//items = append(items, row)
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
		log.Warn().
			Str("methodName", "Mean").
			Str("columnName", name).
			Msg("Column does not exist")
		return 0.0, fmt.Errorf("column %s does not exist", name)
	}

	var kind = d.Columns[name].Kind

	if kind == reflect.String {
		log.Warn().
			Str("methodName", "Mean").
			Str("columnName", name).
			Msg("Column is not numeric")
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
			log.Error().
				Str("methodName", "Mean").
				Str("kind", string(kind)).
				Msg("Unknown type - possible corruption")
			return 0.0, fmt.Errorf("unknown type - possible corruption")
		}
	}

	return avg / float64(d.RowCount), nil
}
