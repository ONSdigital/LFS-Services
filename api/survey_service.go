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
	"sync"
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

func (si SurveyImportHandler) parseGBSurveyFile(tmpfile, datasetName string, week, year, id int) {
	startTime := time.Now()

	si.fileUploads.SetUploadStarted()

	spssData, err := loadSav(tmpfile, types.GBSurveyInput{})
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "parseGBSurveyFile").
			Str("file", datasetName).
			Msg("Cannot import GB SAV file")
		si.fileUploads.SetUploadError(fmt.Sprintf("cannot import GB SAV file %s", err))
		return
	}

	si.Audit.ReferenceDate = time.Now()
	si.Audit.NumObFile = spssData.RowCount
	si.Audit.NumObLoaded = spssData.RowCount
	si.Audit.NumVarFile = spssData.HeaderCount
	si.Audit.NumVarLoaded = spssData.HeaderCount
	si.Audit.FileName = datasetName
	si.Audit.Id = id
	si.Audit.Year = year
	si.Audit.Week = week
	si.Audit.FileSource = types.GBSource

	pipeline := filter.NewGBPipeLine(spssData, &si.Audit)

	columns, data, err := pipeline.RunPipeline()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("preProcessing failed")
		si.fileUploads.SetUploadError(fmt.Sprintf("pre-processing failed %s", err))
		return
	}

	log.Debug().
		Str("datasetName", datasetName).
		Int("numObservationsFile", si.Audit.NumObFile).
		Int("numObservationsLoaded", si.Audit.NumObLoaded).
		Int("numVarFile", si.Audit.NumVarFile).
		Int("numVarLoaded", si.Audit.NumVarLoaded).
		Str("status", "Successful").
		Msg("preProcessing complete")

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		si.fileUploads.SetUploadError(fmt.Sprintf("cannot connect to database %s", err))
		return
	}

	surveyVo := types.SurveyVO{
		Audit:   &si.Audit,
		Records: data,
		Columns: columns,
		Status:  si.fileUploads,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := database.PersistSurvey(surveyVo); err != nil {
			log.Error().
				Err(err).
				Str("datasetName", datasetName).
				Msg("Cannot persist GB survey data")
			si.fileUploads.SetUploadError(fmt.Sprintf("cannot persist GB survey data: %s", err))
			return
		}
	}()

	// todo: add value labels here

	go func() {
		defer wg.Done()
		if err := database.PersistVariableDefinitions(spssData.Header, types.GBSource); err != nil {
			log.Error().
				Err(err).
				Str("datasetName", datasetName).
				Msg("Cannot persist variable definitions (GB)")
			si.fileUploads.SetUploadError(fmt.Sprintf("cannot persist variable definitions (GB): %s", err))
			return
		}
	}()

	wg.Wait()

	log.Debug().
		Str("datasetName", datasetName).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Imported and persisted GB survey data")

	return
}

func (si SurveyImportHandler) parseNISurveyFile(tmpfile, datasetName string, month, year, id int) {
	startTime := time.Now()

	si.fileUploads.SetUploadStarted()

	spssData, err := loadSav(tmpfile, types.NISurveyInput{})
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "parseNISurveyFile").
			Str("file", datasetName).
			Msg("Cannot import NI SAV file")
		si.fileUploads.SetUploadError(fmt.Sprintf("cannot import NI SAV file %s", err))
		return
	}

	// calculate starting week for the month
	weekNo := 1
	for i := 1; i < si.Audit.Month; i++ {
		if i == 2 || i == 5 || i == 8 || i == 11 {
			weekNo = weekNo + 5
		} else {
			weekNo = weekNo + 4
		}
	}

	si.Audit.ReferenceDate = time.Now()
	si.Audit.NumObFile = spssData.RowCount
	si.Audit.NumObLoaded = spssData.RowCount
	si.Audit.NumVarFile = spssData.HeaderCount
	si.Audit.NumVarLoaded = spssData.HeaderCount
	si.Audit.FileName = datasetName
	si.Audit.Id = id
	si.Audit.Year = year
	si.Audit.Month = month
	si.Audit.Week = weekNo
	si.Audit.FileSource = types.NISource

	pipeline := filter.NewNIPipeLine(spssData, &si.Audit)

	columns, data, err := pipeline.RunPipeline()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("preProcessing failed")
		si.fileUploads.SetUploadError(fmt.Sprintf("pre-processing failed: %s", err))
		return
	}

	log.Debug().
		Str("datasetName", datasetName).
		Int("numObservationsFile", si.Audit.NumObFile).
		Int("numObservationsLoaded", si.Audit.NumObLoaded).
		Int("numVarFile", si.Audit.NumVarFile).
		Int("numVarLoaded", si.Audit.NumVarLoaded).
		Str("status", "Successful").
		Msg("preProcessing complete")

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		si.fileUploads.SetUploadError(fmt.Sprintf("cannot connect to database: %s", err))
		return
	}

	surveyVo := types.SurveyVO{
		Audit:   &si.Audit,
		Records: data,
		Columns: columns,
		Status:  si.fileUploads,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := database.PersistSurvey(surveyVo); err != nil {
			log.Error().
				Err(err).
				Str("datasetName", datasetName).
				Msg("Cannot persist NI survey data")
			si.fileUploads.SetUploadError(fmt.Sprintf("cannot persist NI survey data: %s", err))
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := database.PersistVariableDefinitions(spssData.Header, types.NISource); err != nil {
			log.Error().
				Err(err).
				Str("datasetName", datasetName).
				Msg("Cannot persist variable definitions (NI)")
			si.fileUploads.SetUploadError(fmt.Sprintf("cannot persist variable definitions (NI): %s", err))
			return
		}
	}()

	wg.Wait()

	log.Debug().
		Str("datasetName", datasetName).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Imported and persisted NI survey data")

	return
}
