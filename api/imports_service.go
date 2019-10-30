package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/api/filter"
	"services/api/validate"
	"services/dataset"
	"services/db"
	"services/importdata/csv"
	"services/importdata/sav"
	"services/types"
	"services/util"
	"time"
)

func findNIBatch(monthNo, yearNo int) (types.NIBatchItem, error) {

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot connect to database")
		return types.NIBatchItem{}, fmt.Errorf("cannot connect to database: %s", err)
	}

	info, err := database.FindNIBatchInfo(monthNo, yearNo)
	if err != nil {
		log.Error().Err(err)
		return types.NIBatchItem{},
			fmt.Errorf("cannot upload the survey file as the batch for month %d and year %d does not exist", monthNo, yearNo)
	}

	return info, nil
}

func findGBBatch(weekNo, yearNo int) (types.GBBatchItem, error) {

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot connect to database")
		return types.GBBatchItem{}, fmt.Errorf("cannot connect to database: %s", err)
	}

	info, err := database.FindGBBatchInfo(weekNo, yearNo)
	if err != nil {
		log.Error().Err(err)
		return types.GBBatchItem{},
			fmt.Errorf("cannot upload the survey file as the batch for week %d and year %d does not exist",
				weekNo, yearNo)
	}

	return info, nil
}

// TODO: Run in goroutine
func (im ImportsHandler) parseAddressFile(fileName, datasetName string) error {

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

func (im ImportsHandler) parseGBSurveyFile(tmpfile, datasetName string, week, year, id int) error {
	startTime := time.Now()

	log.Debug().
		Str("status", "Successful").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Validator complete")

	records, err := sav.ImportSav(tmpfile)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		log.Warn().
			Str("method", "readSav").
			Msg("The SAV file is empty")
		return fmt.Errorf("the spss file: %s is empty", tmpfile)
	}

	//database, err := db.GetDefaultPersistenceImpl()
	//if err != nil {
	//	log.Error().
	//		Err(err).
	//		Str("datasetName", datasetName).
	//		Msg("Cannot connect to database")
	//	return fmt.Errorf("cannot connect to database: %s", err)
	//}

	if err := populateDatabase(records, types.SurveyInput{}); err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot persist dataset")
		return fmt.Errorf("cannot persist dataset to database: %s", err)
	}

	return nil
}

// TODO: Run in goroutine
func (im ImportsHandler) parseNISurveyFile(tmpfile, datasetName string, month, year, id int) error {
	startTime := time.Now()

	d, err := dataset.NewDataset(datasetName)
	if err != nil {
		return err
	}

	var surveyFilter = filter.NewNISurveyFilter(&d)

	err = d.LoadSav(tmpfile, datasetName, types.SurveyInput{}, surveyFilter)
	if err != nil {
		return err
	}

	startValidation := time.Now()

	val := validate.NewSurveyValidation(&d, validate.GB)
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
