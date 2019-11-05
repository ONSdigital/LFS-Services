package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/api/filter"
	"services/dataset"
	"services/db"
	"services/importdata/sav"
	"services/types"
	"services/util"
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

	rows, err := loadSav(tmpfile, types.GBSurveyInput{})
	if err != nil {
		return err
	}

	headers := rows[0]
	body := rows[1:]

	si.Audit.ReferenceDate = time.Now()
	si.Audit.NumObFile = len(body)
	si.Audit.NumObLoaded = len(body)
	si.Audit.NumVarFile = len(headers)
	si.Audit.NumVarLoaded = len(headers)
	si.Audit.FileName = datasetName
	si.Audit.Id = id
	si.Audit.Year = year
	si.Audit.Week = week
	si.Audit.FileSource = types.GBSource

	pipeline := filter.NewGBPipeLine(rows, &si.Audit)

	columns, data, err := pipeline.RunPipeline()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("pipeline failed")
		return err
	}

	log.Debug().
		Str("datasetName", datasetName).
		Int("numObservationsFile", si.Audit.NumObFile).
		Int("numObservationsLoaded", si.Audit.NumObLoaded).
		Int("numVarFile", si.Audit.NumVarFile).
		Int("numVarLoaded", si.Audit.NumVarLoaded).
		Str("status", "Successful").
		Msg("Pipeline complete")

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		return fmt.Errorf("cannot connect to database: %s", err)
	}

	surveyVo := types.SurveyVO{
		Audit:   &si.Audit,
		Records: data,
		Columns: columns,
	}

	if err := database.PersistSurvey(surveyVo); err != nil {
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

	rows, err := loadSav(tmpfile, types.GBSurveyInput{})
	if err != nil {
		return err
	}

	headers := rows[0]
	body := rows[1:]

	si.Audit.ReferenceDate = time.Now()
	si.Audit.NumObFile = len(body)
	si.Audit.NumObLoaded = len(body)
	si.Audit.NumVarFile = len(headers)
	si.Audit.NumVarLoaded = len(headers)
	si.Audit.FileName = datasetName
	si.Audit.Id = id
	si.Audit.Year = year
	si.Audit.Month = month
	si.Audit.FileSource = types.NISource

	pipeline := filter.NewGBPipeLine(rows, &si.Audit)

	columns, data, err := pipeline.RunPipeline()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("pipeline failed")
		return err
	}

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
		Audit:   &si.Audit,
		Records: data,
		Columns: columns,
	}

	if err := database.PersistSurvey(surveyVo); err != nil {
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
