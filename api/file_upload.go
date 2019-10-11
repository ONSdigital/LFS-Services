package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"services/api/filter"
	"services/api/validate"
	"services/dataset"
	"services/db"
	"time"
)

func (h RestHandlers) fileUpload() error {

	_ = h.r.ParseMultipartForm(64 << 20)

	file, _, err := h.r.FormFile("lfsFile")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error getting formfile")
		return err
	}

	defer func() { _ = file.Close() }()

	fileName := h.r.Form.Get("fileName")
	if fileName == "" {
		log.Error("fileName not set")
		return fmt.Errorf("fileName not set")
	}

	fileType := h.r.Form.Get("fileType")
	if fileType == "" {
		log.Error("fileType not set")
		return fmt.Errorf("fileType not set")
	}

	log.WithFields(log.Fields{
		"fileName": fileName,
		"fileType": fileType,
	}).Debug("Uploading file...")

	startTime := time.Now()

	tmpfile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %s ", err)
	}

	defer func() { _ = os.Remove(tmpfile.Name()) }()

	n, err := io.Copy(tmpfile, file)

	elapsed := time.Now().Sub(startTime)
	log.WithFields(log.Fields{
		"fileName":    fileName,
		"fileType":    fileType,
		"bytesRead":   n,
		"elapsedtime": elapsed,
	}).Debug("File uploaded")

	_ = tmpfile.Close()

	switch fileType {

	case SurveyFile:
		source := h.r.Form.Get("fileSource") // GB or NI
		if source != "GB" && source != "NI" {
			log.Error("fileSource must be NI or GB")
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
		log.WithFields(log.Fields{
			"error":    "filetype not recognised",
			"fileName": fileName,
			"fileType": fileType,
		}).Warn("Error getting formfile")
		return fmt.Errorf("fileType, %s, not recognised", fileType)
	}

	return nil
}

func (h RestHandlers) geogUpload(tmpfile, datasetName string) error {
	return nil
}

func (h RestHandlers) surveyUpload(tmpfile, datasetName, source string) error {
	startTime := time.Now()

	if source == "GB" {

	} else {
		// must be NI
	}
	d, err := dataset.NewDataset(datasetName)
	if err != nil {
		return err
	}

	err = d.LoadSav(tmpfile, datasetName, dataset.Survey{})
	if err != nil {
		return err
	}

	startValidation := time.Now()

	val := validate.NewSurveyValidation(&d)
	validationResponse, err := val.Validate()
	if err != nil {
		log.WithFields(log.Fields{
			"status":       "Failed",
			"errorMessage": err,
			"elapsedTime":  time.Now().Sub(startValidation),
		}).Warn("Validator complete")
		return err
	}

	if validationResponse.ValidationResult == validate.ValidationFailed {
		log.WithFields(log.Fields{
			"status":       "Failed",
			"errorMessage": validationResponse.ErrorMessage,
			"elapsedTime":  time.Now().Sub(startValidation),
		}).Warn("Validator complete")
		return fmt.Errorf(validationResponse.ErrorMessage)
	}

	log.WithFields(log.Fields{
		"status":      "Successful",
		"elapsedTime": time.Now().Sub(startValidation),
	}).Debug("Validator complete")

	f := filter.NewSurveyFilter(&d)

	f.DropColumns()
	f.RenameColumns()

	err = f.AddVariables()
	if err != nil {
		log.WithFields(log.Fields{
			"datasetName":  datasetName,
			"errorMessage": err.Error(),
		}).Error(err)
		return err
	}

	log.WithFields(log.Fields{
		"status": "Successful",
	}).Debug("Filtering complete")

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.WithFields(log.Fields{
			"datasetName":  datasetName,
			"errorMessage": err.Error(),
		}).Error("cannot connect to database")
		return fmt.Errorf("cannot connect to database: %s", err)
	}

	if err := database.PersistDataset(d); err != nil {
		log.WithFields(log.Fields{
			"datasetName":  datasetName,
			"errorMessage": err.Error(),
		}).Error("cannot persist dataset to database")
		return fmt.Errorf("cannot persist dataset to database: %s", err)
	}

	log.WithFields(log.Fields{
		"datasetName": datasetName,
		"rowCount":    d.NumRows(),
		"columnCount": d.ColumnCount,
		"elapsedTime": time.Now().Sub(startTime),
	}).Debug("Imported and persisted dataset")

	return nil
}
