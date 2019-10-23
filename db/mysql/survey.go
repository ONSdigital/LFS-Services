package mysql

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/config"
	"services/dataset"
	"services/types"
	"services/util"
	"strconv"
	"strings"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

var columnsTable string

func init() {
	columnsTable = config.Config.Database.ColumnsTable
	if columnsTable == "" {
		panic("columns table configuration not set")
	}
}

func (s MySQL) DeleteColumnData(name string) error {
	col := s.DB.Collection(columnsTable)
	res := col.Find("table_name", name)
	if res == nil {
		return nil
	}
	if err := res.Delete(); err != nil {
		return err
	}
	return nil
}

func (s MySQL) insertColumnData(tx sqlbuilder.Tx, columns types.Columns) error {

	col := tx.Collection(columnsTable)
	_, err := col.Insert(columns)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) UnpersistSurveyDataset(tableName string) (dataset.Dataset, error) {
	d, err := dataset.NewDataset(tableName)

	startTime := time.Now()
	log.Info().Msg("starting unpersist")

	if err != nil {
		log.Error().
			Err(err).
			Str("methodName", "UnpersistSurveyDataset").
			Msg("Cannot create a new DataSet")
		return dataset.Dataset{}, fmt.Errorf("cannot create a new DataSet: %s", err)
	}

	log.Info().Msg("starting unpersist into Dataset")

	req := s.DB.Collection(columnsTable).Find().Where("table_name = '" + tableName + "'").OrderBy("column_number")
	var column types.Columns
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
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int8 - possible corruption")
				}
			case reflect.Int, reflect.Int32, reflect.Uint32:
				s[i], err = strconv.ParseInt(v, 10, 32)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int32 - possible corruption")
				}
			case reflect.Int64, reflect.Uint64:
				s[i], err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int64 - possible corruption")
				}
			case reflect.Float32:
				s[i], err = strconv.ParseFloat(v, 32)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on float32 - possible corruption")
				}
			case reflect.Float64:
				s[i], err = strconv.ParseFloat(v, 64)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on float64 - possible corruption")
				}
			default:
				log.Error().
					Err(err).
					Str("methodName", "UnpersistSurveyDataset").
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
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Data unpersisted")
	return d, nil
}

func (s MySQL) PersistSurveyDataset(d dataset.Dataset) error {
	var kBuffer bytes.Buffer

	startTime := time.Now()
	log.Debug().Msg("Starting persistence into DB")

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
					Str("methodName", "PersistSurveyDataset").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption")
				return fmt.Errorf("unknown type - possible corruption")
			}

			if i != len(column.Rows)-1 {
				kBuffer.WriteString(",")
			}
			i++
		}

		column := types.Columns{
			TableName:    d.DatasetName,
			ColumnName:   colName,
			ColumnNumber: column.ColNo,
			Kind:         int(columnKind),
			Rows:         kBuffer.String(),
		}

		if err := s.insertColumnData(tx, column); err != nil {
			return fmt.Errorf("cannot insert column, error: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	var f = DBAudit{s}

	if err := f.AuditFileUploadEvent(d); err != nil {
		log.Error().
			Err(err).
			Msg("AuditFileUpload failed")
		return fmt.Errorf("AuditFileUpload, error: %s", err)
	}

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Survey data persisted")

	return nil
}
