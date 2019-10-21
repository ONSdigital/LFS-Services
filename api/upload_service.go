package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"services/api/filter"
	"services/api/validate"
	"services/dataset"
	"services/db"
	"services/importdata/csv"
	"services/types"
	"services/util"
	"time"
)

func (h RestHandlers) fileUpload(fileType string) error {

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
		Str("elapsedTime", util.FmtDuration(startTime)).
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

	case AddressFile:
		if err := h.addressUpload(tmpfile.Name(), fileName); err != nil {
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

func (h RestHandlers) uploadAddress(w http.ResponseWriter, r *http.Request) Response {

	if err := h.fileUpload(AddressFile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrorResponse{Status: Error, ErrorMessage: err.Error()}
	}

	return OkayResponse{OK}
}

func (h RestHandlers) addressUpload(fileName, datasetName string) error {

	startTime := time.Now()

	rows, err := csv.ImportCSVToSlice(fileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "addressUpload").
			Str("file", fileName).
			Msg("Cannot import CSV file")
		return fmt.Errorf("cannot import CSV file %w", err)
	}

	if len(rows) < 2 {
		log.Warn().
			Str("method", "addressUpload").
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

func (h RestHandlers) uploadSurvey(w http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	runId := vars["runId"]

	if runId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn().Msg("runId not set")
		return ErrorResponse{Status: Error, ErrorMessage: "runId not set"}
	}

	if err := h.fileUpload(SurveyFile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrorResponse{Status: Error, ErrorMessage: err.Error()}
	}

	return OkayResponse{OK}
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
		surveyFilter = filter.NewNISurveyFilter(&d)
	}

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
