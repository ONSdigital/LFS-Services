package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"os"
	"services/api/filter"
	"services/api/validate"
	"services/dataset"
	"services/db"
	"services/types"
	"time"
)

func (h RestHandlers) fileUpload() error {

	_ = h.r.ParseMultipartForm(64 << 20)

	file, _, err := h.r.FormFile("lfsFile")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error getting formfile")
		return err
	}

	defer func() { _ = file.Close() }()

	fileName := h.r.Form.Get("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		return fmt.Errorf("fileName not set")
	}

	fileType := h.r.Form.Get("fileType")
	if fileType == "" {
		log.Error().Msg("fileName not set")
		return fmt.Errorf("fileType not set")
	}

	log.Debug().
		Str("fileName", fileName).
		Str("fileType", fileType).
		Msg("Uploading file")

	startTime := time.Now()

	tmpfile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %s ", err)
	}

	defer func() { _ = os.Remove(tmpfile.Name()) }()

	n, err := io.Copy(tmpfile, file)

	log.Debug().
		Str("fileName", fileName).
		Str("fileType", fileType).
		Int64("bytesRead", n).
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("File uploaded")

	_ = tmpfile.Close()

	switch fileType {

	case SurveyFile:
		source := h.r.Form.Get("fileSource") // GB or NI
		if source != "GB" && source != "NI" {
			log.Error().Msg("fileSource must be NI or GB")
			return fmt.Errorf("invalid fileSource or fileSource not set - must be GB or NI")
		}
		if err := h.surveyUpload(tmpfile.Name(), fileName, source); err != nil {
			return err
		}

	case GeogFile:
		if err := h.geogUpload(tmpfile.Name(), fileName); err != nil {
			return err
		}

	default:
		log.Warn().
			Str("fileName", fileName).
			Str("fileType", fileType).
			Int64("bytesRead", n).
			Str("error", "filetype not recognised").
			Msg("Error getting formfile")
		return fmt.Errorf("fileType, %s, not recognised", fileType)
	}

	return nil
}

func (h RestHandlers) geogUpload(tmpfile, datasetName string) error {
	return nil
}

func (h RestHandlers) surveyUpload(tmpfile, datasetName, source string) error {
	startTime := time.Now()

	d, err := dataset.NewDataset(datasetName)
	if err != nil {
		return err
	}

	var surveyFilter types.Filter
	if source == "GB" {
		surveyFilter = filter.NewGBSurveyFilter(&d)
	} else {
		// not GB so must be NI
		surveyFilter = filter.NewNISurveyFilter(&d)
	}

	err = d.LoadSav(tmpfile, datasetName, dataset.Survey{}, surveyFilter)
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
			TimeDiff("elapsedTime", time.Now(), startValidation).
			Msg("Validator failed")
		return err
	}

	if validationResponse.ValidationResult == validate.ValidationFailed {
		log.Warn().
			Str("status", "Failed").
			Str("errorMessage", validationResponse.ErrorMessage).
			TimeDiff("elapsedTime", time.Now(), startValidation).
			Msg("Validator failed")
		return fmt.Errorf(validationResponse.ErrorMessage)
	}

	log.Debug().
		Str("status", "Successful").
		TimeDiff("elapsedTime", time.Now(), startValidation).
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

	if err := database.PersistDataset(d); err != nil {
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
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("Imported and persisted dataset")

	return nil
}
