package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/api/filter"
	"services/api/validate"
	"services/dataset"
	"services/db"
	"services/importdata/csv"
	"services/types"
	"services/util"
	"time"
)

// TODO: Run in goroutine
func (h RestHandlers) parseAddressFile(fileName, datasetName string) error {

	startTime := time.Now()

	rows, err := csv.ImportCSVToSlice(fileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "parseAddressFile").
			Str("file", fileName).
			Msg("Cannot import CSV file")
		return fmt.Errorf("cannot import CSV file %w", err)
	}

	if len(rows) < 2 {
		log.Warn().
			Str("method", "parseAddressFile").
			Msg("The CSV file is empty")
		return fmt.Errorf("csv file: %s is empty", fileName)
	}

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		return fmt.Errorf("cannot connect to database: %s", err)
	}

	if err := database.PersistAddressDataset(rows[0], rows[1:]); err != nil {
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

/*
NI and GB are in seperate functions for now - need to see how different they really are
*/

// TODO: Run in goroutine
func (h RestHandlers) parseGBSurveyFile(tmpfile, datasetName string) error {
	startTime := time.Now()

	d, err := dataset.NewDataset(datasetName)
	if err != nil {
		return err
	}

	var surveyFilter = filter.NewGBSurveyFilter(&d)

	err = d.LoadSav(tmpfile, datasetName, types.Survey{}, surveyFilter)
	if err != nil {
		return err
	}

	startValidation := time.Now()

	val := validate.NewSurveyValidation(&d)
	validationResponse, err := val.Validate()
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

	cnt, err := surveyFilter.AddVariables()
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

	if err := database.PersistSurveyDataset(d); err != nil {
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

// TODO: Run in goroutine
func (h RestHandlers) parseNISurveyFile(tmpfile, datasetName string) error {
	startTime := time.Now()

	d, err := dataset.NewDataset(datasetName)
	if err != nil {
		return err
	}

	var surveyFilter = filter.NewNISurveyFilter(&d)

	err = d.LoadSav(tmpfile, datasetName, types.Survey{}, surveyFilter)
	if err != nil {
		return err
	}

	startValidation := time.Now()

	val := validate.NewSurveyValidation(&d)
	validationResponse, err := val.Validate()
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

	cnt, err := surveyFilter.AddVariables()
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

	if err := database.PersistSurveyDataset(d); err != nil {
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
