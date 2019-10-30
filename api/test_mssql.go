package api

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/denisenkom/go-mssqldb"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/util"
	"strconv"
	"strings"
	"time"
)

func populateDatabase(rows [][]string, out interface{}) error {

	var err error
	var columnCount = 0
	type Column struct {
		ColNo int
		Kind  reflect.Kind
	}

	startTime := time.Now()

	headers := rows[0]
	body := rows[1:]

	var columns = make(map[string]Column, len(headers))

	addColumn := func(name string, columnType reflect.Kind) error {
		if _, ok := columns[name]; ok {
			log.Warn().
				Str("method", "AddColumn").
				Str("column", name).
				Msg("Column already exists")
			return fmt.Errorf("column %s already exists", name)
		}

		col := Column{}
		col.Kind = columnType
		columnCount++
		columns[name] = col
		return nil
	}

	logStructError := func(methodName, variableName string, kind reflect.Kind, newType string) {
		log.Error().
			Str("methodName", methodName).
			Str("variable", variableName).
			Str("convertFrom", string(kind)).
			Str("convertTo", newType).
			Msg("Camnnot convert type")
	}

	t1 := reflect.TypeOf(out)

	//var f = filter.NewGBSurveyFilter(nil)
	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		// skip columns that are marked as being dropped
		//if f.DropColumn(strings.ToUpper(a.Name)) {
		//	continue
		//}
		if a.Name == "Union" {
			a.Name = "UnionCol"
		}
		if err := addColumn(a.Name, a.Type.Kind()); err != nil {
			log.Error().
				Err(err).
				Str("methodName", "populateDataset").
				Str("columnName", strings.ToUpper(a.Name)).
				Str("columnType", string(a.Type.Kind())).
				Msg("Cannot create column")
			return fmt.Errorf("cannot create column: %w", err)
		}
	}

	//for a := range headers {
	//	headers[a] = strings.ToUpper(headers[a])
	//}

	log.Info().Msg("Connecting...")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
		"localhost", "lfs", "lfs", 1433)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return err
	}

	txn, err := db.Begin()
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msg("Start bulk load")

	_, err = txn.Exec("create table #staging(cs nvarchar(max))")
	if err != nil {
		log.Fatal().Err(err)
		return err
	}

	stmt, err := txn.Prepare(mssql.CopyIn("#staging", mssql.BulkOptions{
		CheckConstraints:  false,
		FireTriggers:      false,
		KeepNulls:         false,
		KilobytesPerBatch: 0,
		RowsPerBatch:      5000,
		Order:             nil,
		Tablock:           false,
	}, "cs"))
	if err != nil {
		log.Fatal().Err(err)
	}

	defer func() {
		_ = db.Close()
	}()

	for _, spssRow := range body {

		var row bytes.Buffer

		for j := 0; j < len(spssRow); j++ {
			if len(spssRow) != len(headers) {
				log.Error().
					Err(err).
					Str("methodName", "populateDatabase").
					Str("rowSize", string(len(spssRow))).
					Str("columnSize", string(len(headers))).
					Msg("Header is out of alignment with row")
				return fmt.Errorf("header is out of alignment with row")
			}

			//header := strings.ToUpper(headers[j])
			header := headers[j]
			// extract the tagged columns only
			if _, ok := columns[headers[j]]; !ok {
				continue
			}

			a := spssRow[j]
			if a == "" || a == "NaN" {
				a = "NULL"
			}

			kind := columns[headers[j]].Kind

			switch kind {

			case reflect.String:
				if a == "NULL" {
					break
				}

				row.WriteString(fmt.Sprintf("<%s>", header))
				by := []byte(a)
				err = xml.EscapeText(&row, by)
				if err != nil {
					logStructError("populateDataset", header, kind, "string")
					return fmt.Errorf("cannot convert %s into an escaped xml string", a)
				}
				row.WriteString(fmt.Sprintf("</%s>", header))

			case reflect.Int8, reflect.Uint8:
				if a == "NULL" {
					break
				}
				i, err := strconv.ParseInt(a, 0, 8)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int8")
					return fmt.Errorf("cannot convert %s into an Int8", a)
				}
				row.WriteString(fmt.Sprintf("<%s>%d</%s>", header, i, header))

			case reflect.Int, reflect.Int32, reflect.Uint32:
				if a == "NULL" {
					break
				}
				i, err := strconv.ParseInt(a, 0, 32)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int32")
					return fmt.Errorf("cannot convert %s into an Int32", a)
				}
				row.WriteString(fmt.Sprintf("<%s>%d</%s>", header, i, header))

			case reflect.Int64, reflect.Uint64:
				if a == "NULL" {
					break
				}
				i, err := strconv.ParseInt(a, 0, 64)
				if err != nil {
					logStructError("populateDataset", header, kind, "Int64")
					return fmt.Errorf("cannot convert %s into an Int64", a)
				}
				row.WriteString(fmt.Sprintf("<%s>%d</%s>", header, i, header))

			case reflect.Float32:
				if a == "NULL" {
					break
				}
				i, err := strconv.ParseFloat(a, 32)

				if err != nil {
					logStructError("populateDataset", header, kind, "Float32")
					return fmt.Errorf("cannot convert %s into an Float32", a)
				}
				row.WriteString(fmt.Sprintf("<%s>%f</%s>", header, i, header))

			case reflect.Float64:
				if a == "NULL" {
					break
				}
				i, err := strconv.ParseFloat(a, 64)
				if err != nil {
					logStructError("populateDataset", header, kind, "Float64")
					return fmt.Errorf("cannot convert %s into an Float64", a)
				}
				row.WriteString(fmt.Sprintf("<%s>%f</%s>", header, i, header))

			default:
				logStructError("populateDataset", header, kind, "Unknown")
				return fmt.Errorf("cannot convert struct variable type from SPSS type")
			}
		}

		// call the skipRow filter
		//if f.SkipRow(row) {
		//	continue
		//}

		_, err = stmt.Exec(row.String())
		if err != nil {
			log.Error().Err(err)
			panic(err)
		}

	}

	//m := make(map[string]Column, d.NumColumns())
	//
	//for k, v := range d.Columns {
	//	to, ok := filter.RenameColumns(k)
	//	if ok {
	//		m[to] = v
	//	} else {
	//		m[k] = v
	//	}
	//}
	//
	//d.Columns = m

	result, err := stmt.Exec()
	if err != nil {
		log.Fatal().Err(err)
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal().Err(err)
		return err
	}

	log.Debug().Msg("Start copy from staging")

	_, err = txn.Exec("insert into Survey(cs) select cs from #staging")
	if err != nil {
		log.Fatal().Err(err).Msg("copy from temp table to Survey failed")
		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal().Err(err)
		return err
	}

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("End copy from staging")

	rowCount, _ := result.RowsAffected()

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startTime)).
		Int64("rowsCopied", rowCount).
		Msg("Imported and persisted dataset")

	return nil
}
