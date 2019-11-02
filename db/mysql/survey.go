package mysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"reflect"
	"services/config"
	"services/dataset"
	"services/types"
	"services/util"
	"strconv"
	"strings"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var surveyTable string

func init() {
	surveyTable = config.Config.Database.SurveyTable
	if surveyTable == "" {
		panic("survey table configuration not set")
	}
}

func (s MySQL) DeleteSurveyData(name string) (bool, error) {
	//cnt, err := s.DB.Collection(surveyTable).Find("file_name", name).Count()
	cnt, err := s.DB.Collection(surveyTable).Find(db.Cond{"file_name": name}).Count()

	if err != nil {
		return false, err
	}
	if cnt == 0 {
		return false, nil
	}
	q := s.DB.DeleteFrom(surveyTable).Where("file_name", name)
	_, err = q.Exec()
	if err != nil {
		return true, err
	}
	return true, nil
}

func (s MySQL) insertSurveyData(tx sqlbuilder.Tx, survey types.SurveyRow) error {

	col := tx.Collection(surveyTable)
	_, err := col.Insert(survey)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) UnpersistSurveyDataset(tableName string) (dataset.Dataset, error) {
	//d, err := dataset.NewDataset(tableName)
	//
	//startTime := time.Now()
	//log.Info().Msg("starting unpersist")
	//
	//if err != nil {
	//	log.Error().
	//		Err(err).
	//		Str("methodName", "UnpersistSurveyDataset").
	//		Msg("Cannot create a new DataSet")
	//	return dataset.Dataset{}, fmt.Errorf("cannot create a new DataSet: %s", err)
	//}
	//
	//log.Info().Msg("starting unpersist into Dataset")
	//
	//req := s.DB.Collection(surveyTable).Find().
	//	Where("file_name = '" + tableName + "'").
	//	OrderBy("column_number")
	//
	//var column types.SurveyRow
	//for req.Next(&column) {
	//	a := strings.Split(column.Rows, ",")
	//	s := make([]interface{}, len(a))
	//	for i, v := range a {
	//		switch reflect.Kind(column.Kind) {
	//		case reflect.String:
	//			s[i] = v
	//		case reflect.Int8, reflect.Uint8:
	//			s[i], err = strconv.ParseInt(v, 10, 8)
	//			if err != nil {
	//				return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int8 - possible corruption")
	//			}
	//		case reflect.Int, reflect.Int32, reflect.Uint32:
	//			s[i], err = strconv.ParseInt(v, 10, 32)
	//			if err != nil {
	//				return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int32 - possible corruption")
	//			}
	//		case reflect.Int64, reflect.Uint64:
	//			s[i], err = strconv.ParseInt(v, 10, 64)
	//			if err != nil {
	//				return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int64 - possible corruption")
	//			}
	//		case reflect.Float32:
	//			s[i], err = strconv.ParseFloat(v, 32)
	//			if err != nil {
	//				return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on float32 - possible corruption")
	//			}
	//		case reflect.Float64:
	//			s[i], err = strconv.ParseFloat(v, 64)
	//			if err != nil {
	//				return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on float64 - possible corruption")
	//			}
	//		default:
	//			log.Error().
	//				Err(err).
	//				Str("methodName", "UnpersistSurveyDataset").
	//				Str("type", string(reflect.Kind(column.Kind))).
	//				Msg("Unknown type - possible corruption")
	//			return dataset.Dataset{}, fmt.Errorf("unknown type - possible corruption")
	//		}
	//	}
	//	col := dataset.Column{
	//		ColNo: column.ColumnNumber,
	//		Kind:  reflect.Kind(column.Kind),
	//		Rows:  s,
	//	}
	//	d.Columns[column.ColumnName] = col
	//	d.ColumnCount++
	//	d.RowCount = len(s)
	//}
	//
	//log.Debug().
	//	Err(err).
	//	Str("elapsedTime", util.FmtDuration(startTime)).
	//	Msg("Data unpersisted")
	return dataset.Dataset{}, nil
}

func (s MySQL) PersistSurveyDataset(d dataset.Dataset, vo types.SurveyVO) error {
	var kBuffer bytes.Buffer

	startTime := time.Now()
	log.Debug().Msg("Starting persistence into DB")

	found, err := s.DeleteSurveyData(vo.FileName)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Delete existing survey data failed")
		return fmt.Errorf("delete existing survey data failed, error: %s", err)
	}

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	type record struct {
		Name    []interface{} `json:"name"`
		Records []interface{} `json:"value"`
	}

	for _, column := range d.Columns {
		kBuffer.Reset()

		records := record{Records: make([]interface{}, d.RowCount)}

		var i = 0
		columnKind := column.Kind
		for k, v := range column.Rows {
			switch columnKind {
			case reflect.String:
				if v == "NULL" {
					v = nil
				}
				records.Records[k] = v

			case reflect.Int8, reflect.Uint8:
				records.Records[k] = v

			case reflect.Int, reflect.Int32, reflect.Uint32:
				records.Records[k] = v

			case reflect.Int64, reflect.Uint64:
				records.Records[k] = v

			case reflect.Float32:
				num := v.(float64)
				if math.IsNaN(num) {
					records.Records[k] = nil
				} else {
					records.Records[k] = v
				}

			case reflect.Float64:
				num := v.(float64)
				if math.IsNaN(num) {
					records.Records[k] = nil
				} else {
					records.Records[k] = v
				}

			default:
				log.Error().
					Str("methodName", "PersistSurveyDataset").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption")
				return fmt.Errorf("unknown type - possible corruption")
			}

			i++
		}

		_, err := json.Marshal(records)
		if err != nil {
			return fmt.Errorf("cannot marshal json, error: %s", err)
		}
		row := types.SurveyRow{
			Id:         vo.Id,
			FileName:   vo.FileName,
			FileSource: vo.FileSource,
			Week:       vo.Week,
			Month:      vo.Month,
			Year:       vo.Year,
		}

		if err := s.insertSurveyData(tx, row); err != nil {
			return fmt.Errorf("cannot insert survey row, error: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	vo.NumObLoaded = d.NumObLoaded
	vo.NumObFile = d.NumObFile
	vo.NumVarLoaded = d.NumVarLoaded
	vo.NumVarFile = d.NumVarFile

	weekOrMonth := func() int {
		if vo.FileSource == types.GBSource {
			return vo.Week
		}
		return vo.Month
	}()

	var updateMonthlyStatus func(int, int, int) error
	if vo.FileSource == types.GBSource {
		updateMonthlyStatus = s.updateGBBatch
	} else {
		updateMonthlyStatus = s.updateNIBatch
	}

	if found {
		if err := updateMonthlyStatus(weekOrMonth, vo.Year, types.FileReloaded); err != nil {
			return err
		}
		if err := s.insertAudit(vo, types.FileReloaded, "File re-uploaded successfully"); err != nil {
			return err
		}
	} else {
		if err := updateMonthlyStatus(weekOrMonth, vo.Year, types.FileUploaded); err != nil {
			return err
		}
		if err := s.insertAudit(vo, types.FileUploaded, "File uploaded successfully"); err != nil {
			return err
		}
	}

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("GBSurveyInput data persisted")

	return nil
}

func (s MySQL) insertAudit(vo types.SurveyVO, status int, message string) error {
	audit := types.Audit{
		Id:            vo.Id,
		FileName:      surveyTable,
		FileSource:    vo.FileSource,
		Week:          vo.Week,
		Month:         vo.Month,
		Year:          vo.Year,
		ReferenceDate: time.Now(),
		NumVarFile:    vo.NumVarFile,
		NumVarLoaded:  vo.NumVarLoaded,
		NumObFile:     vo.NumObFile,
		NumObLoaded:   vo.NumObLoaded,
		Status:        status,
		Message:       message,
	}

	var f = DBAudit{s}
	if err := f.AuditFileUploadEvent(audit); err != nil {
		log.Error().
			Err(err).
			Msg("AuditFileUpload failed")
		return fmt.Errorf("AuditFileUpload, error: %s", err)
	}
	return nil
}

func (s MySQL) PersistSurvey(rows [][]string, vo types.SurveyVO, filter types.Filter) error {

	log.Debug().Msg("Starting persistence into DB")

	_, err := s.DeleteSurveyData(vo.FileName)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Delete existing survey data failed")
		return fmt.Errorf("delete existing survey data failed, error: %s", err)
	}

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	type Column struct {
		Name  string
		Skip  bool
		ColNo int
		Kind  reflect.Kind
	}

	headers := rows[0]
	body := rows[1:]

	out := types.GBSurveyInput{}

	columns := make([]Column, len(headers))

	t1 := reflect.TypeOf(out)

	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		col := Column{}
		// skip columns that are marked as being dropped
		if filter.DropColumn(strings.ToUpper(a.Name)) {
			col.Skip = true
			continue
		}
		col.Skip = false
		col.Kind = a.Type.Kind()
		col.Name = a.Name
		col.ColNo = i
		columns[i] = col
	}

	var kBuffer bytes.Buffer

	for _, v := range body {
		kBuffer.Reset()
		kBuffer.WriteString("{")

		for col_no, val := range v {

			if columns[col_no].Skip {
				continue
			}

			kBuffer.WriteString("\"" + columns[col_no].Name + "\":")

			columnKind := columns[col_no].Kind
			switch columnKind {
			case reflect.String:
				if val == "NULL" || val == "" {
					val = "null"
					kBuffer.WriteString(val)
				} else {
					kBuffer.WriteString("\"")
					kBuffer.WriteString(jsonEscape(val))
					kBuffer.WriteString("\"")
				}

			case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
				kBuffer.WriteString(val)

			case reflect.Float32, reflect.Float64:
				if val == "" || val == "NULL" {
					val = "0.0"
				}
				f, err := strconv.ParseFloat(val, 64)
				if err != nil {
					log.Error().
						Str("methodName", "PersistSurveyDataset").
						Int("type", int(columnKind)).
						Msg("field is not a float")
					return fmt.Errorf("field is not a float")
				}
				if math.IsNaN(f) {
					kBuffer.WriteString("null")
				} else {
					kBuffer.WriteString(val)
				}

			default:
				log.Error().
					Str("methodName", "PersistSurveyDataset").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption or structure does not map to file")
				return fmt.Errorf("unknown type - possible corruption or structure does not map to file")
			}

			if col_no != len(v)-1 {
				kBuffer.WriteString(",")
			} else {
				kBuffer.WriteString("}")
			}
		}

		row := types.SurveyRow{
			Id:         vo.Id,
			FileName:   vo.FileName,
			FileSource: vo.FileSource,
			Week:       vo.Week,
			Month:      vo.Month,
			Year:       vo.Year,
			Columns:    kBuffer.String(),
		}

		if err := s.insertSurveyDataV1(tx, row); err != nil {
			return fmt.Errorf("cannot insert survey row, error: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	return nil
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func (s MySQL) insertSurveyDataV1(tx sqlbuilder.Tx, survey types.SurveyRow) error {

	col := tx.Collection(surveyTable)
	_, err := col.Insert(survey)
	if err != nil {
		return err
	}

	return nil
}
