package redis

import "C"
import (
	"bytes"
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"pds-go/lfs/importdata/sav"
	"pds-go/lfs/io/spss"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var globalLock = sync.Mutex{}

type Dataset struct {
	tableName     string
	TableMetaData map[string]reflect.Kind
	pool          *redis.Pool
	mux           *sync.Mutex
	logger        *log.Logger
	count         int
}

const BulkSize = 5000
const COMMA = ","

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
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)

	return nil
}

func NewDataset(name string, logger *log.Logger) (*Dataset, error) {

	pool := newPool()
	conn := pool.Get()
	err := ping(conn)

	if err != nil {
		return nil, fmt.Errorf(" -> NewDataset: cannot connect to redis, error: %s", err)
	}

	_, _ = conn.Do("DEL", name)

	mux := sync.Mutex{}
	return &Dataset{name, nil, pool, &mux, logger, 0}, nil
}

type fromFileFunc func(fileName string, out interface{}) (dataset Dataset, err error)

func (d *Dataset) logLoad(from fromFileFunc) fromFileFunc {
	return func(fileName string, out interface{}) (dataset Dataset, err error) {
		startTime := time.Now()
		res, err := from(fileName, out)
		a := time.Now().Sub(startTime)
		d.logger.Printf("file load processed in %s", a)
		return res, err
	}
}

func (d *Dataset) FromSav(fileName string, out interface{}) (dataset Dataset, err error) {
	return d.logLoad(d.readSav)(fileName, out)
}

func (d *Dataset) readSav(in string, out interface{}) (dataset Dataset, err error) {
	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		return Dataset{}, fmt.Errorf(" -> FromSav: %T is not a struct type", out)
	}

	start := time.Now()

	records, err := sav.ImportSav(in)
	if err != nil {
		return Dataset{}, err
	}

	if len(records) == 0 {
		return Dataset{}, fmt.Errorf(" -> createDataset: spss file: %s is empty", in)
	}

	elapsed := time.Since(start)
	d.logger.Info(fmt.Sprintf("read sav file (%d records) in %s", len(records), elapsed))

	start = time.Now()
	i, er := d.createDataset(in, records, out)
	if er != nil {
		return Dataset{}, er
	}
	elapsed = time.Since(start)
	d.logger.Info(fmt.Sprintf("created dataset (%d records) in %s", len(records), elapsed))

	return i, nil
}

func (d Dataset) AddColumn(name string, columnType spss.ColumnTypes) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	//sqlStmt := fmt.Sprintf("alter table %s add %s %s", d.tableName, name, columnType)
	//_, err := d.DB.Exec(sqlStmt)
	//if err != nil {
	//	return fmt.Errorf(" -> AddColumn: cannot insert column: %s", err)
	//}
	return nil
}

func (d *Dataset) createDataset(fileName string, rows [][]string, out interface{}) (Dataset, error) {
	_, file := filepath.Split(fileName)
	var extension = filepath.Ext(file)
	var name = file[0 : len(file)-len(extension)]
	d, er := NewDataset(name, d.logger)

	if er != nil {
		return Dataset{}, fmt.Errorf(" -> createDataset: cannot create a new DataSet: %s", er)
	}

	d.logger.Println("starting import into Dataset")

	t1 := reflect.TypeOf(out)
	d.TableMetaData = make(map[string]reflect.Kind)

	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		d.TableMetaData[a.Name] = a.Type.Kind()
	}

	headers := rows[0]
	body := rows[1:]
	count := 0

	rs := make([]map[string]interface{}, 0)

	var wg sync.WaitGroup
	errorChannel := make(chan error)

	//v := math.Ceil(float64(len(body) / BulkSize))
	v := len(body) / BulkSize
	wg.Add(v)

	for _, spssRow := range body {
		row := make(map[string]interface{})

		for j := 0; j < len(spssRow); j++ {
			if len(spssRow) != len(headers) {
				return Dataset{}, fmt.Errorf(" -> createDataset: header is out of alignment with row. row size: %d, column size: %d", len(spssRow), len(headers))
			}
			header := headers[j]
			// extract the columns we are interested in
			if _, ok := d.TableMetaData[headers[j]]; !ok {
				continue
			}

			// check type is valid
			a := spssRow[j]
			if a == "" {
				a = "NULL"
			}

			kind := d.TableMetaData[headers[j]]
			switch kind {

			case reflect.String:
				break
			case reflect.Int8, reflect.Uint8:
				i, err := strconv.ParseInt(a, 0, 8)
				if err != nil {
					return Dataset{}, fmt.Errorf(" -> createDataset: cannot convert %s into an Int8", a)
				}
				row[header] = i

			case reflect.Int, reflect.Int32, reflect.Uint32:
				i, err := strconv.ParseInt(a, 0, 32)
				if err != nil {
					return Dataset{}, fmt.Errorf(" -> createDataset: cannot convert %s into an Int32", a)
				}
				row[header] = i

			case reflect.Int64, reflect.Uint64:
				i, err := strconv.ParseInt(a, 0, 64)
				if err != nil {
					return Dataset{}, fmt.Errorf(" -> createDataset: cannot convert %s into an Int64", a)
				}
				row[header] = i

			case reflect.Float32:
				i, err := strconv.ParseFloat(a, 32)
				if err != nil {
					return Dataset{}, fmt.Errorf(" -> createDataset: cannot convert %s into an Float32", a)
				}
				row[header] = i

			case reflect.Float64:
				i, err := strconv.ParseFloat(a, 64)
				if err != nil {
					return Dataset{}, fmt.Errorf(" -> createDataset: cannot convert %s into an Float64", a)
				}
				row[header] = i

			default:
				return Dataset{}, fmt.Errorf(" -> createDataset: cannot convert struct variable type from SPSS type")
			}

			row[header] = spssRow[j]
		}

		rs = append(rs, row)
		count++

		if count%BulkSize == 0 {
			go func(rs []map[string]interface{}) {
				defer wg.Done()
				if err := d.BulkInsert(rs); err != nil {
					errorChannel <- fmt.Errorf(" -> createDataset: cannot create row: %s", err)
					//return Dataset{}, fmt.Errorf(" -> createDataset: cannot create row: %s", err)
				}
			}(rs)
			rs = nil
		}
	}

	wg.Wait()

	select {
	case err := <-errorChannel:
		return Dataset{}, err
	default:
	}

	if rs != nil {
		if err := d.BulkInsert(rs); err != nil {
			return Dataset{}, fmt.Errorf(" -> createDataset: cannot create row: %s", err)
		}
	}

	return *d, nil
}

func (d *Dataset) BulkInsert(values []map[string]interface{}) (err error) {
	var kBuffer bytes.Buffer

	conn := d.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	for _, row := range values {
		kBuffer.Reset()
		kBuffer.WriteString("{")
		d.count++
		rowLabel := fmt.Sprintf("%s:%d", d.tableName, d.count)

		var i = 0
		for k, v := range row {
			kBuffer.WriteString(fmt.Sprintf("\"%s\":", k))
			if d.TableMetaData[k] == reflect.String {
				a := fmt.Sprintf("%s", v)
				a = strings.Replace(a, "'", `''`, -1)
				kBuffer.WriteString("\"" + a + "\"")
			} else {
				kBuffer.WriteString(fmt.Sprintf("%s", v))
			}
			if i != len(row)-1 {
				kBuffer.WriteString(",")
			} else {
				kBuffer.WriteString("}")
			}
			i++
		}

		_, err := conn.Do("SET", rowLabel, kBuffer.String())
		if err != nil {
			panic(err)
		}

	}

	return
}
