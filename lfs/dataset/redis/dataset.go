package redis

import "C"
import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
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
	colNo int
	kind  reflect.Kind
	rows  []interface{}
}

type Dataset struct {
	tableName   string
	Columns     map[string]Column
	pool        *redis.Pool
	mux         *sync.Mutex
	logger      *log.Logger
	rowCount    int
	columnCount int
}

const (
	BulkSize              = 5000
	InitialRowCapacity    = 20000
	InitialColumnCapacity = 2000
	COMMA                 = ","
)

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 5,
		// max number of connections
		MaxActive: 100,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	if _, err := redis.String(pong, err); err != nil {
		return err
	}

	return nil
}

func NewDataset(name string, logger *log.Logger) (Dataset, error) {

	pool := newPool()
	conn := pool.Get()
	err := ping(conn)

	if err != nil {
		return Dataset{}, fmt.Errorf(" -> NewDataset: cannot connect to redis, error: %s", err)
	}

	_, _ = conn.Do("DEL", name)

	mux := sync.Mutex{}
	cols := make(map[string]Column, InitialColumnCapacity)
	return Dataset{name, cols, pool, &mux, logger, 0, 0}, nil
}

type fromFileFunc func(fileName string, out interface{}) error

func (d *Dataset) logLoad(from fromFileFunc) fromFileFunc {
	return func(fileName string, out interface{}) error {
		startTime := time.Now()
		err := from(fileName, out)
		a := time.Now().Sub(startTime)
		d.logger.Printf("file load processed in %s", a)
		return err
	}
}

func (d *Dataset) LoadSav(fileName string, out interface{}) error {
	return d.logLoad(d.readSav)(fileName, out)
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
	d.logger.Info(fmt.Sprintf("created dataset (%d records) in %s", d.rowCount, elapsed))

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
		col.rows = append(col.rows, v)
		d.Columns[k] = col
	}
	d.rowCount++
	return nil
}

func (d *Dataset) AddColumn(name string, columnType reflect.Kind) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if _, ok := d.Columns[name]; ok {
		return fmt.Errorf(" -> AddColumn: Column %s already exists", name)
	}

	col := Column{}
	col.kind = columnType
	col.colNo = d.columnCount
	col.rows = make([]interface{}, 0, InitialRowCapacity)
	d.Columns[name] = col
	d.columnCount++

	if d.rowCount == 0 {
		return nil
	}

	// Add empty rows if we hae existing data
	for i := 0; i < d.rowCount; i++ {
		switch columnType {
		case reflect.String:
			col.rows = append(col.rows, "")
		case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
			col.rows = append(col.rows, 0)
		case reflect.Float32, reflect.Float64:
			col.rows = append(col.rows, 0.0)
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

	a := d.orderedColumns()
	m := make(map[string]Column, InitialColumnCapacity)

	var colNo = 0
	for _, v := range a {
		if v != name {
			var col Column
			col.rows = m[v].rows
			col.kind = m[v].kind
			col.colNo = colNo
			m[v] = col
			colNo++
		}
	}

	d.Columns = m
	d.columnCount--
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

	d.logger.Println("starting import into Dataset")

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

			kind := d.Columns[headers[j]].kind
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

func (d Dataset) orderedColumns() []string {
	var keys = make([]string, d.columnCount)
	for k, v := range d.Columns {
		keys[v.colNo] = k
	}
	return keys
}

func (d *Dataset) getByRow(maxRows int, maxCols int) ([]string, [][]string) {
	cnt := 0
	var header []string
	var items [][]string

	if maxCols > d.columnCount {
		maxCols = d.columnCount
	}

	for _, v := range d.orderedColumns() {
		if cnt > maxCols-1 {
			break
		}
		header = append(header, v)
		cnt++
	}

	if maxRows > d.rowCount {
		maxRows = d.rowCount
	}

	// for each header, get MaxRows
	for j := 0; j < maxRows; j++ {
		var row []string
		for _, b := range header {
			r := d.Columns[b].rows[j]
			kind := d.Columns[b].kind

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
	return d.columnCount
}

func (d Dataset) NumRows() int {
	return d.rowCount
}

func (d Dataset) Mean(name string) (float64, error) {
	if _, ok := d.Columns[name]; !ok {
		return 0.0, fmt.Errorf(" -> Mean: Column %s does not exist", name)
	}

	var kind = d.Columns[name].kind

	if kind == reflect.String {
		return 0.0, errors.New(fmt.Sprintf(" -> Mean: column %s is not numeric", name))
	}

	var avg = 0.0

	for _, v := range d.Columns[name].rows {
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

	return avg / float64(d.rowCount), nil
}

//func (d *Dataset) BulkInsert(values []map[string]interface{}) (err error) {
//	var kBuffer bytes.Buffer
//
//	conn := d.pool.Get()
//	defer func() {
//		_ = conn.Close()
//	}()
//
//	for _, row := range values {
//		kBuffer.Reset()
//		kBuffer.WriteString("{")
//		d.count++
//		rowLabel := fmt.Sprintf("%s:%d", d.tableName, d.count)
//
//		var i = 0
//		for k, v := range row {
//			kBuffer.WriteString(fmt.Sprintf("\"%s\":", k))
//			if d.TableMetaData[k] == reflect.String {
//				a := fmt.Sprintf("%s", v)
//				a = strings.Replace(a, "'", `''`, -1)
//				kBuffer.WriteString("\"" + a + "\"")
//			} else {
//				kBuffer.WriteString(fmt.Sprintf("%s", v))
//			}
//			if i != len(row)-1 {
//				kBuffer.WriteString(",")
//			} else {
//				kBuffer.WriteString("}")
//			}
//			i++
//		}
//
//		_, err := conn.Do("SET", rowLabel, kBuffer.String())
//		if err != nil {
//			panic(err)
//		}
//
//	}
//
//	return
//}
