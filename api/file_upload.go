package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"services/api/filters"
	"services/dataset"
	"services/db"
	"time"
)

func (h RestHandlers) fileUpload() error {

	_ = h.r.ParseMultipartForm(32 << 20)

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

	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

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
		if err := h.surveyUpload(tmpfile.Name(), fileName); err != nil {
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
		}).Error("Error getting formfile")
		return fmt.Errorf("fileType, %s, not recognised", fileType)
	}

	return nil
}

func (h RestHandlers) geogUpload(tmpfile, datasetName string) error {
	return nil
}

func (h RestHandlers) surveyUpload(tmpfile, datasetName string) error {
	startTime := time.Now()
	logger := log.New()

	d, err := dataset.NewDataset(datasetName, logger)
	if err != nil {
		return err
	}

	err = d.LoadSav(tmpfile, datasetName, dataset.Survey{})
	if err != nil {
		return err
	}

	filter := filters.NewSurveyFilter(h.log)
	err = filter.Validate()
	if err != nil {
		return fmt.Errorf("validation failed: %s", err)
	}

	filter.DropColumns(&d)
	filter.RenameColumns(&d)

	err = db.GetDefaultPersistenceImpl(logger).PersistDataset(d)
	if err != nil {
		return fmt.Errorf("GetDefaultPersistenceImpl failed: %s", err)
	}

	elapsed := time.Now().Sub(startTime)
	log.WithFields(log.Fields{
		"datasetName": datasetName,
		"rowCount":    d.NumRows(),
		"columnCount": d.ColumnCount,
		"elapsedTime": elapsed,
	}).Info("imported and persisted dataset")

	return nil
}
