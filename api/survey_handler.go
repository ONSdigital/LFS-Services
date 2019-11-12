package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"services/api/ws"
	"services/types"
	"strconv"
	"sync"
	"time"
)

type SurveyImportHandler struct {
	Audit            types.Audit
	fileUploads      *types.WSMessage
	uploadInProgress bool // we can only handle a single survey upload at a time
	mutux            *sync.Mutex
}

func NewSurveyHandler() *SurveyImportHandler {
	return &SurveyImportHandler{
		Audit:            types.Audit{},
		fileUploads:      nil,
		uploadInProgress: false,
		mutux:            &sync.Mutex{},
	}
}

func (si *SurveyImportHandler) setUpload(val bool) {
	si.mutux.Lock()
	defer si.mutux.Unlock()
	si.uploadInProgress = val
}

/*
Upload GB survey file.
The mutux lock is a little awkward to ensure it is locked as soon as possible to avoid race conditions
*/
func (si *SurveyImportHandler) SurveyUploadGBHandler(w http.ResponseWriter, r *http.Request) {

	if si.uploadInProgress {
		log.Error().Msg("Survey file is currently being uploaded")
		ErrorResponse{Status: Error, ErrorMessage: "survey file is currently being uploaded"}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	si.setUpload(true)

	vars := mux.Vars(r)
	week := vars["week"]
	year := vars["year"]

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("File name not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	weekNo, err := strconv.Atoi(week)
	if err != nil {
		log.Error().Msg("Week is not an integer")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		log.Error().Msg("Year is not an integer")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	gbInfo, err := FindGBBatch(weekNo, yearNo)
	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	a := ws.NewFileUploads()
	si.fileUploads = a.Add(fileName)

	go func() {
		defer func() {
			si.setUpload(false)
			_ = os.Remove(tmpfile)
		}()
		si.parseGBSurveyFile(tmpfile, fileName, weekNo, yearNo, gbInfo.Id)
	}()

	InProgressResponse{OK, time.Now().String(), "request submitted"}.sendResponse(w, r)
}

func (si *SurveyImportHandler) SurveyUploadNIHandler(w http.ResponseWriter, r *http.Request) {

	if si.uploadInProgress {
		log.Error().Msg("Survey file is currently being uploaded")
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: "survey file is currently being uploaded"}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	si.setUpload(true)

	vars := mux.Vars(r)
	month := vars["month"]
	year := vars["year"]

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("File name not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	monthNo, err := strconv.Atoi(month)
	if err != nil {
		log.Error().Msg("Month is not an integer")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		log.Error().Msg("Year is not an integer")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	niInfo, err := FindNIBatch(monthNo, yearNo)
	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		si.setUpload(false)
		return
	}

	a := ws.NewFileUploads()
	si.fileUploads = a.Add(fileName)

	go func() {
		defer func() {
			si.setUpload(false)
			_ = os.Remove(tmpfile)
		}()
		si.parseNISurveyFile(tmpfile, fileName, monthNo, yearNo, niInfo.Id)
	}()

	InProgressResponse{OK, time.Now().String(), "request submitted"}.sendResponse(w, r)

}
