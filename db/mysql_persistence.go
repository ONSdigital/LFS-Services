package db

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/config"
	"services/dataset"
	"services/types"
	"strconv"
	"strings"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type MySQL struct {
	DB sqlbuilder.Database
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

	log.Debug().
		Str("databaseName", config.Config.Database.Database).
		Msg("Connecting to database")

	sess, err := mysql.Open(settings)

	if err != nil {
		log.Error().
			Err(err).
			Str("databaseName", config.Config.Database.Database).
			Msg("Cannot connect to database")
		return err
	}

	log.Debug().
		Str("databaseName", config.Config.Database.Database).
		Msg("Connected to database")

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
	}

	log.Debug().
		Int("MaxPoolSize", poolSize).
		Int("MaxIdleConnections", maxIdle).
		Dur("MaxLifetime", maxLifetime*time.Second).
		Msg("Connection Attributes")

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
	d, err := dataset.NewDataset(tableName)

	startTime := time.Now()
	log.Info().Msg("starting unpersist")

	if err != nil {
		log.Error().
			Err(err).
			Str("methodName", "UnpersistDataset").
			Msg("Cannot create a new DataSet")
		return dataset.Dataset{}, fmt.Errorf("cannot create a new DataSet: %s", err)
	}

	log.Info().Msg("starting unpersist into Dataset")

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
				log.Error().
					Err(err).
					Str("methodName", "UnpersistDataset").
					Str("type", string(reflect.Kind(column.Kind))).
					Msg("Unknown type - possible corruption")
				return dataset.Dataset{}, fmt.Errorf("unknown type - possible corruption")
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

	log.Debug().
		Err(err).
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("Data unpersisted")
	return d, nil
}

func (s MySQL) PersistDataset(d dataset.Dataset) error {
	var kBuffer bytes.Buffer

	startTime := time.Now()
	log.Debug().
		Str("tableName", d.DatasetName).
		Msg("Starting persistence into DB")

	_ = s.DeleteColumnData(d.DatasetName)

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
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
				log.Error().
					Str("methodName", "PersistDataset").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption")
				return fmt.Errorf("unknown type - possible corruption")
			}

			if i != len(column.Rows)-1 {
				kBuffer.WriteString(",")
			}
			i++
		}

		if err := s.insertColumnData(tx, d.DatasetName, colName, column.ColNo, int(columnKind), kBuffer.String()); err != nil {
			return fmt.Errorf("cannot insert column, error: %s", err)
		}
	}

	if err := s.auditFileUpload(tx, d); err != nil {
		log.Error().
			Err(err).
			Msg("AuditFileUpload failed")
		return fmt.Errorf("AuditFileUpload, error: %s", err)
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	log.Debug().
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("Data persisted")

	return nil
}

type Audit struct {
	ReferenceDate time.Time `db:"reference_date"`
	FileName      string    `db:"file_name"`
	NumVarFile    int       `db:"num_var_file"`
	NumVarLoaded  int       `db:"num_var_loaded"`
	NumObFile     int       `db:"num_ob_file"`
	NumObLoaded   int       `db:"num_ob_loaded"`
}

func (s MySQL) auditFileUpload(tx sqlbuilder.Tx, d dataset.Dataset) error {
	a := Audit{
		FileName:      d.DatasetName,
		ReferenceDate: time.Now(),
		NumVarFile:    d.NumVarFile,
		NumVarLoaded:  d.NumVarLoaded,
		NumObFile:     d.NumObFile,
		NumObLoaded:   d.NumObLoaded,
	}
	dbAudit := tx.Collection("upload_audit")
	_, err := dbAudit.Insert(a)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) GetUserID(user string) (types.UserCredentials, error) {
	var creds types.UserCredentials

	col := s.DB.Collection("users")
	res := col.Find("username", user)

	if res == nil {
		return creds, fmt.Errorf("user %s not found", user)
	}

	defer func() { _ = res.Close() }()

	ok := res.Next(&creds)
	if !ok {
		return creds, fmt.Errorf("user %s not found", user)
	}
	return creds, nil
}
