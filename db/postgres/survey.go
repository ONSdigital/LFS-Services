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

func (s Postgres) DeleteSurveyData(audit types.Audit) (bool, error) {
	col := s.DB.Collection(surveyTable)

	res := col.Find(db.Cond{
		"file_name": audit.FileName, "week": audit.Week,
		"year": audit.Year, "file_source": audit.FileSource,
	})

	defer func() { _ = res.Close() }()

	if res.Err() != nil {
		return false, res.Err()
	}

	count, err := res.Count()
	if err != nil {
		return true, err
	}
	if count == 0 {
		return false, nil
	}
	err = res.Delete()

	if err != nil {
		return true, err
	}
	return true, nil
}

func (s Postgres) PersistSurvey(vo types.SurveyVO) error {

	log.Debug().Msg("Starting persistence into DB")

	existingDataset, err := s.DeleteSurveyData(*vo.Audit)
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

	for cnt, v := range body {
		var rowMap = make(map[string]*interface{})

		for colNo, val := range v {

			if columns[colNo].Skip {
				continue
			}

			columnKind := columns[colNo].Kind
			switch columnKind {
			case reflect.String:
				if val == "NULL" || val == "" {
					continue
					//rowMap[columns[colNo].Name] = nil
				} else {
					var ms interface{} = val
					rowMap[columns[colNo].Name] = &ms
				}

			case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
				i64, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Error().
						Str("methodName", "PersistSurvey").
						Int("type", int(columnKind)).
						Msg("field is not an int")
					return fmt.Errorf("field is not an int")
				}
				var ms interface{} = i64
				if i64 == 0 {
					continue
					//ms = nil
				}
				rowMap[columns[colNo].Name] = &ms

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
					//rowMap[columns[colNo].Name] = nil
					continue
				} else {
					var ms interface{} = f
					rowMap[columns[colNo].Name] = &ms
				}

			default:
				log.Error().
					Str("methodName", "PersistSurvey").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption or structure does not map to file")
				return fmt.Errorf("unknown type - possible corruption or structure does not map to file")
			}

			if colNo == len(v)-1 {
				var perc = (float64(cnt) / float64(len(body))) * 100
				vo.Status.SetPercentage(perc)
			}
		}

		re, err := json.Marshal(rowMap)
		if err != nil {
			return fmt.Errorf("json marshall failed: %s", err)
		}

		row := types.SurveyRow{
			Id:         vo.Audit.Id,
			FileName:   vo.Audit.FileName,
			FileSource: vo.Audit.FileSource,
			Week:       vo.Audit.Week,
			Month:      vo.Audit.Month,
			Year:       vo.Audit.Year,
			Columns:    string(re),
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

func (s Postgres) insertSurveyData(tx sqlbuilder.Tx, survey types.SurveyRow) error {

	col := tx.Collection(surveyTable)
	_, err := col.Insert(survey)
	if err != nil {
		return err
	}

	return nil
}
