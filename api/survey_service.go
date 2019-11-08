package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/api/filter"
	"services/db"
	"services/importdata/sav"
	"services/types"
	"services/util"
	"time"
)

func loadSav(in string, out interface{}) (types.SavImportData, error) {

	// ensure out is a struct
	if reflect.ValueOf(out).Kind() != reflect.Struct {
		log.Error().
			Str("method", "readSav").
			Msg("The output interface is not a struct")
		return types.SavImportData{}, fmt.Errorf("%T is not a struct type", out)
	}

	start := time.Now()

	records, err := sav.ImportSav(in)
	if err != nil {
		return types.SavImportData{}, err
	}

	if records.RowCount == 0 {
		log.Warn().
			Str("method", "readSav").
			Msg("The SAV file is empty")
		return types.SavImportData{}, fmt.Errorf("the spss file: %s is empty", in)
	}

	log.Debug().
		Str("file", in).
		Int("records", records.RowCount-1).
		Str("elapsedTime", util.FmtDuration(start)).
		Msg("Read Sav file")

	return records, nil
}

func (si SurveyImportHandler) parseGBSurveyFile(tmpfile, datasetName string, week, year, id int) error {
	startTime := time.Now()

	spssData, err := loadSav(tmpfile, types.GBSurveyInput{})
	if err != nil {
		return err
	}

	headers, body := sav.SPSSDatatoArray(spssData)

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

	pipeline := filter.NewGBPipeLine(headers, body, &si.Audit)

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

	spssData, err := loadSav(tmpfile, types.GBSurveyInput{})
	if err != nil {
		return err
	}

	headers, body := sav.SPSSDatatoArray(spssData)

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

	pipeline := filter.NewNIPipeLine(headers, body, &si.Audit)

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
