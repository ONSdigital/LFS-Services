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
	"strconv"
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

	return nil
}

func (s MySQL) insertAudit(vo types.SurveyVO, status int, message string) error {
	audit := types.Audit{
		Id:            vo.Audit.Id,
		FileName:      surveyTable,
		FileSource:    vo.Audit.FileSource,
		Week:          vo.Audit.Week,
		Month:         vo.Audit.Month,
		Year:          vo.Audit.Year,
		ReferenceDate: time.Now(),
		NumVarFile:    vo.Audit.NumVarFile,
		NumVarLoaded:  vo.Audit.NumVarLoaded,
		NumObFile:     vo.Audit.NumObFile,
		NumObLoaded:   vo.Audit.NumObLoaded,
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

func (s MySQL) PersistSurvey(vo types.SurveyVO) error {

	log.Debug().Msg("Starting persistence into DB")

	_, err := s.DeleteSurveyData(vo.Audit.FileName)
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

	body := vo.Records[1:]

	columns := vo.Columns

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
			Id:         vo.Audit.Id,
			FileName:   vo.Audit.FileName,
			FileSource: vo.Audit.FileSource,
			Week:       vo.Audit.Id,
			Month:      vo.Audit.Month,
			Year:       vo.Audit.Year,
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
