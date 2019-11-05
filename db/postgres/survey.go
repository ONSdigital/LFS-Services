package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"reflect"
	"services/config"
	"services/types"
	"strconv"
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

func (s Postgres) DeleteSurveyData(name string) (bool, error) {
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

func (s Postgres) PersistSurvey(vo types.SurveyVO) error {

	log.Debug().Msg("Starting persistence into DB")

	existingDataset, err := s.DeleteSurveyData(vo.Audit.FileName)
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
						Str("methodName", "PersistSurvey").
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
					Str("methodName", "PersistSurvey").
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
			Week:       vo.Audit.Week,
			Month:      vo.Audit.Month,
			Year:       vo.Audit.Year,
			Columns:    kBuffer.String(),
		}

		if err := s.insertSurveyData(tx, row); err != nil {
			_ = tx.Rollback()
			vo.Audit.Status = types.UploadError
			vo.Audit.Message = "Insert survey row failed"
			_ = s.AuditFileUploadEvent(*vo.Audit)
			return fmt.Errorf("cannot insert survey row, error: %s", err)
		}
	}

	if existingDataset {
		vo.Audit.Status = types.FileReloaded
	} else {
		vo.Audit.Status = types.UploadFinished
	}
	vo.Audit.Message = "File Uploaded"
	if err := s.AuditFileUploadEvent(*vo.Audit); err != nil {
		log.Error().
			Err(err).
			Msg("Audit event failed")
		_ = tx.Rollback()
		return fmt.Errorf("audit event failed, error: %s", err)
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

func (s Postgres) insertSurveyData(tx sqlbuilder.Tx, survey types.SurveyRow) error {

	col := tx.Collection(surveyTable)
	_, err := col.Insert(survey)
	if err != nil {
		return err
	}

	return nil
}
