package postgres

import (
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

	defer vo.Status.SetUploadFinished()

	type ColumnDetails struct {
		Name  *string     `json:"name"`
		Value interface{} `json:"value"`
		Label *string     `json:"label"`
	}

	var details = make([]ColumnDetails, 0)

	for _, v := range body {
		for colNo, val := range v {

			if columns[colNo].Skip {
				continue
			}

			a := ColumnDetails{
				Name:  &columns[colNo].Name,
				Value: nil,
				Label: &columns[colNo].Label,
			}

			if columns[colNo].Label == "" {
				a.Label = nil
			}

			columnKind := columns[colNo].Kind
			switch columnKind {
			case reflect.String:
				if val == "NULL" || val == "" {
					a.Value = nil
				} else {
					a.Value = val
				}

			case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
				a.Value, err = strconv.Atoi(val)
				if err != nil {
					log.Error().
						Str("methodName", "PersistSurvey").
						Int("type", int(columnKind)).
						Msg("field is not an int")
					return fmt.Errorf("field is not an int")
				}

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
					a.Value = nil
				} else {
					a.Value = f
				}

			default:
				log.Error().
					Str("methodName", "PersistSurvey").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption or structure does not map to file")
				return fmt.Errorf("unknown type - possible corruption or structure does not map to file")
			}

			details = append(details, a)
		}

		a, err := json.Marshal(details)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Cannot marshal survey columns structure")
			return fmt.Errorf("cannot marshal survey columns structure, error: %s", err)
		}

		row := types.SurveyRow{
			Id:         vo.Audit.Id,
			FileName:   vo.Audit.FileName,
			FileSource: vo.Audit.FileSource,
			Week:       vo.Audit.Week,
			Month:      vo.Audit.Month,
			Year:       vo.Audit.Year,
			Columns:    string(a),
		}

		if err := s.insertSurveyData(tx, row); err != nil {
			_ = tx.Rollback()
			vo.Audit.Status = types.UploadError
			vo.Audit.Message = "Insert survey row failed"
			_ = s.AuditFileUploadEvent(*vo.Audit)
			log.Error().
				Err(err).
				Int("week", row.Week).
				Int("month", row.Month).
				Int("year", row.Year).
				Msg("Cannot insert survey row")
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
