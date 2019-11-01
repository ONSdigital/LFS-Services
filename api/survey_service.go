package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"reflect"
	"services/api/filter"
	"services/api/validate"
	"services/dataset"
	"services/db"
	"services/importdata/sav"
	"services/types"
	"services/util"
	"strconv"
	"strings"
	"time"
)

func loadSav(in string, out interface{}) ([][]string, error) {

	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		log.Error().
			Str("method", "readSav").
			Msg("The output interface is not a struct")
		return nil, fmt.Errorf("%T is not a struct type", out)
	}

	start := time.Now()

	records, err := sav.ImportSav(in)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		log.Warn().
			Str("method", "readSav").
			Msg("The CSV file is empty")
		return nil, fmt.Errorf("the spss file: %s is empty", in)
	}

	log.Debug().
		Str("file", in).
		Int("records", len(records)-1).
		Str("elapsedTime", util.FmtDuration(start)).
		Msg("Read Sav file")

	return records, nil
}

func (si SurveyImportHandler) parseGBSurveyFile(tmpfile, datasetName string, week, year, id int) error {
	startTime := time.Now()

	rows, err := loadSav(tmpfile, types.SurveyInput{})
	if err != nil {
		return err
	}

	startValidation := time.Now()

	headers := rows[0]
	body := rows[1:]

	si.Audit.ReferenceDate = time.Now()
	si.Audit.NumObFile = len(body)
	si.Audit.NumObLoaded = len(body)
	si.Audit.NumVarFile = len(headers)
	si.Audit.NumVarLoaded = len(headers)

	val := validate.NewSurveyValidation(validate.GB, &headers, &body)
	validationResponse, err := val.Validate(week, year)
	if err != nil {
		log.Warn().
			Err(err).
			Str("status", "Failed").
			Str("elapsedTime", util.FmtDuration(startValidation)).
			Msg("Validator failed")
		return err
	}

	if validationResponse.ValidationResult == validate.ValidationFailed {
		log.Warn().
			Str("status", "Failed").
			Str("errorMessage", validationResponse.ErrorMessage).
			Str("elapsedTime", util.FmtDuration(startTime)).
			Msg("Validator failed")
		return fmt.Errorf(validationResponse.ErrorMessage)
	}

	log.Debug().
		Str("status", "Successful").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Validator complete")

	var surveyFilter = filter.NewGBSurveyFilter(&si.Audit)

	cnt, err := surveyFilter.AddVariables(&headers, &body)
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName)
		return err
	}

	//// add the number of variables added
	si.Audit.NumVarLoaded = si.Audit.NumVarLoaded + cnt

	log.Debug().
		Str("datasetName", datasetName).
		Int("numObservationsFile", si.Audit.NumObFile).
		Int("numObservationsLoaded", si.Audit.NumObLoaded).
		Int("numVarFile", si.Audit.NumVarFile).
		Int("numVarLoaded", si.Audit.NumVarLoaded).
		Str("status", "Successful").
		Msg("Filtering complete")

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		return fmt.Errorf("cannot connect to database: %s", err)
	}

	surveyVo := types.SurveyVO{
		Id:         id,
		FileName:   datasetName,
		FileSource: "GB",
		Week:       week,
		Month:      0,
		Year:       year,
	}

	if err := database.PersistSurvey(rows, surveyVo, surveyFilter); err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot persist dataset")
		return fmt.Errorf("cannot persist dataset to database: %s", err)
	}

	log.Debug().
		Str("datasetName", datasetName).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Imported and persisted dataset")

	return nil
}

func (si SurveyImportHandler) parseNISurveyFile(tmpfile, datasetName string, month, year, id int) error {
	startTime := time.Now()

	d, err := dataset.NewDataset(datasetName)
	if err != nil {
		return err
	}

	var surveyFilter = filter.NewNISurveyFilter(&si.Audit)

	rows, err := loadSav(tmpfile, types.SurveyInput{})
	if err != nil {
		return err
	}

	startValidation := time.Now()

	headers := rows[0]
	body := rows[1:]

	val := validate.NewSurveyValidation(validate.NI, &headers, &body)
	validationResponse, err := val.Validate(month, year)
	if err != nil {
		log.Warn().
			Err(err).
			Str("status", "Failed").
			Str("elapsedTime", util.FmtDuration(startValidation)).
			Msg("Validator failed")
		return err
	}

	if validationResponse.ValidationResult == validate.ValidationFailed {
		log.Warn().
			Str("status", "Failed").
			Str("errorMessage", validationResponse.ErrorMessage).
			Str("elapsedTime", util.FmtDuration(startTime)).
			Msg("Validator failed")
		return fmt.Errorf(validationResponse.ErrorMessage)
	}

	log.Debug().
		Str("status", "Successful").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Validator complete")

	cnt, err := surveyFilter.AddVariables(&headers, &body)
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName)
		return err
	}

	// add the number of variables added
	d.NumVarLoaded = d.NumVarLoaded + cnt

	log.Debug().
		Str("datasetName", datasetName).
		Int("numObservationsFile", d.NumObFile).
		Int("numObservationsLoaded", d.NumObLoaded).
		Int("numVarFile", d.NumVarFile).
		Int("numVarLoaded", d.NumVarLoaded).
		Str("status", "Successful").
		Msg("Filtering complete")

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		return fmt.Errorf("cannot connect to database: %s", err)
	}

	surveyVo := types.SurveyVO{
		Id:         id,
		FileName:   d.DatasetName,
		FileSource: "NI",
		Week:       0,
		Month:      month,
		Year:       year,
	}

	if err := database.PersistSurveyDataset(d, surveyVo); err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot persist dataset")
		return fmt.Errorf("cannot persist dataset to database: %s", err)
	}

	log.Debug().
		Str("datasetName", datasetName).
		Int("rowCount", d.NumRows()).
		Int("columnCount", d.NumColumns()).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Imported and persisted dataset")

	return nil
}

type Column struct {
	Name  string
	Skip  bool
	ColNo int
	Kind  reflect.Kind
}

func getSurveyStructure(rows [][]string, vo types.SurveyVO, filter types.Filter) ([]types.SurveyRow, error) {
	headers := rows[0]
	body := rows[1:]

	out := types.SurveyInput{}
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

	surveyRows := make([]types.SurveyRow, len(body))

	for rowNo, v := range body {
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
					return nil, fmt.Errorf("field is not a float")
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
				return nil, fmt.Errorf("unknown type - possible corruption or structure does not map to file")
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

		surveyRows[rowNo] = row
	}

	return surveyRows, nil
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}
