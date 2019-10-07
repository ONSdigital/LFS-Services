package db

import (
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"reflect"
	"services/config"
	"services/dataset"
	"strconv"
	"strings"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type MySQL struct {
	DB  sqlbuilder.Database
	log *logger.Logger
}

var columnsTable string

func init() {
	columnsTable = config.Config.Database.ColumnsTable
	if columnsTable == "" {
		panic("columns table configuration not set")
	}
}

func (s *MySQL) Connect() error {

	var settings = mysql.ConnectionURL{
		Database: config.Config.Database.Database,
		Host:     config.Config.Database.Server,
		User:     config.Config.Database.User,
		Password: config.Config.Database.Password,
	}

	s.log.Debug("Connecting to database: ", config.Config.Database.Database)
	sess, err := mysql.Open(settings)

	if err != nil {
		s.log.Fatal(fmt.Errorf("cannot open database connection %v", err))
		return err
	}

	s.log.Debug(fmt.Sprintf("Connected to database: %s", config.Config.Database.Database))

	if config.Config.Database.Verbose {
		sess.SetLogging(true)
	}

	s.DB = sess

	poolSize := config.Config.Database.ConnectionPool.MaxPoolSize
	maxIdle := config.Config.Database.ConnectionPool.MaxIdleConnections
	maxLifetime := config.Config.Database.ConnectionPool.MaxLifetimeSeconds

	if maxLifetime > 0 {
		maxLifetime = maxLifetime * time.Second
		sess.SetConnMaxLifetime(maxLifetime)
		s.log.Debug("MaxLifetime: ", maxLifetime)
	}

	s.log.Debug("MaxPoolSize: ", poolSize)
	s.log.Debug("MaxIdleConnections: ", maxIdle)

	sess.SetMaxOpenConns(poolSize)
	sess.SetMaxIdleConns(maxIdle)

	return nil
}

func (s MySQL) DeleteColumnData(name string) error {
	col := s.DB.Collection("columns")
	res := col.Find("table_name", name)
	if res == nil {
		return nil
	}
	if err := res.Delete(); err != nil {
		return err
	}
	return nil
}

func (s MySQL) Close() {
	if s.DB != nil {
		_ = s.DB.Close()
	}
}

type Column struct {
	TableName    string `db:"table_name"`
	ColumnName   string `db:"column_name"`
	ColumnNumber int    `db:"column_number"`
	Kind         int    `db:"kind"`
	Rows         string `db:"rows"`
}

func (s MySQL) insertColumnData(tx sqlbuilder.Tx, tableName string, columnName string, columnNumber int, kind int, rows string) error {

	c := Column{tableName, columnName, columnNumber, kind, rows}
	col := tx.Collection(columnsTable)
	_, err := col.Insert(c)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) UnpersistDataset(tableName string) (dataset.Dataset, error) {
	d, err := dataset.NewDataset(tableName, s.log)

	startTime := time.Now()
	s.log.Info("starting unpersist")

	if err != nil {
		return dataset.Dataset{}, fmt.Errorf(" -> populateDataset: cannot create a new DataSet: %s", err)
	}

	s.log.Info("starting unpersist into Dataset")
	if err := s.Connect(); err != nil {
		return dataset.Dataset{}, fmt.Errorf(" -> PersistData: cannot connect to database, error: %s", err)
	}

	defer s.Close()

	req := s.DB.Collection(columnsTable).Find().Where("table_name = '" + tableName + "'").OrderBy("column_number")
	var column Column
	for req.Next(&column) {
		a := strings.Split(column.Rows, ",")
		s := make([]interface{}, len(a))
		for i, v := range a {
			switch reflect.Kind(column.Kind) {
			case reflect.String:
				s[i] = v
			case reflect.Int8, reflect.Uint8:
				s[i], err = strconv.ParseInt(v, 10, 8)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistDataset: unpersist error on int8 - possible corruption")
				}
			case reflect.Int, reflect.Int32, reflect.Uint32:
				s[i], err = strconv.ParseInt(v, 10, 32)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistDataset: unpersist error on int32 - possible corruption")
				}
			case reflect.Int64, reflect.Uint64:
				s[i], err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistDataset: unpersist error on int64 - possible corruption")
				}
			case reflect.Float32:
				s[i], err = strconv.ParseFloat(v, 32)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistDataset: unpersist error on float32 - possible corruption")
				}
			case reflect.Float64:
				s[i], err = strconv.ParseFloat(v, 64)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistDataset: unpersist error on float64 - possible corruption")
				}
			default:
				panic(fmt.Errorf(" -> getByRow: unknown type - possible corruption"))
			}
		}
		col := dataset.Column{
			ColNo: column.ColumnNumber,
			Kind:  reflect.Kind(column.Kind),
			Rows:  s,
		}
		d.Columns[column.ColumnName] = col
		d.ColumnCount++
		d.RowCount = len(s)
	}

	a := time.Now().Sub(startTime)
	s.log.Info("data unpersisted in ", a.String())
	return d, nil
}

func (s MySQL) PersistDataset(d dataset.Dataset) error {
	var kBuffer bytes.Buffer

	startTime := time.Now()
	s.log.Info("starting persistence")

	if err := s.Connect(); err != nil {
		return fmt.Errorf(" -> PersistData: cannot connect to database, error: %s", err)
	}

	defer s.Close()

	_ = s.DeleteColumnData(d.TableName)

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		return fmt.Errorf(" -> PersistData: cannot start a transaction, error: %s", err)
	}

	for colName, column := range d.Columns {
		kBuffer.Reset()

		var i = 0
		columnKind := column.Kind
		for _, v := range column.Rows {
			switch columnKind {
			case reflect.String:
				a := fmt.Sprintf("%s", v)
				a = strings.Replace(a, "'", `''`, -1)
				//kBuffer.WriteString("\"" + a + "\"")
				kBuffer.WriteString(a)
			case reflect.Int8, reflect.Uint8:
				kBuffer.WriteString(fmt.Sprintf("%d", v))
			case reflect.Int, reflect.Int32, reflect.Uint32:
				kBuffer.WriteString(fmt.Sprintf("%d", v))
			case reflect.Int64, reflect.Uint64:
				kBuffer.WriteString(fmt.Sprintf("%d", v))
			case reflect.Float32:
				kBuffer.WriteString(fmt.Sprintf("%f", v))
			case reflect.Float64:
				kBuffer.WriteString(fmt.Sprintf("%g", v))
			default:
				panic(fmt.Errorf(" -> getByRow: unknown type - possible corruption"))
			}

			if i != len(column.Rows)-1 {
				kBuffer.WriteString(",")
			}
			i++
		}

		if err := s.insertColumnData(tx, d.TableName, colName, column.ColNo, int(columnKind), kBuffer.String()); err != nil {
			return fmt.Errorf(" -> PersistData: cannot insert column, error: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(" -> PersistData: commit failed, error: %s", err)
	}

	a := time.Now().Sub(startTime)
	s.log.Info(fmt.Sprintf("data persisted in %s", a.String()))

	return nil
}
