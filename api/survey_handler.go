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

func (si *SurveyImportHandler) SurveyUploadGBHandler(w http.ResponseWriter, r *http.Request) {

	si.mutux.Lock()

	w.Header().Set("Content-Type", "application/json")

	if si.uploadInProgress {
		log.Error().Msg("survey file is currently being uploaded")
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: "survey file is currently being uploaded"}.sendResponse(w, r)
		si.mutux.Unlock()
		return
	}
	si.uploadInProgress = true
	si.mutux.Unlock()

	vars := mux.Vars(r)
	week := vars["week"]
	year := vars["year"]

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("week", week).
		Str("yeay", year).
		Msg("Received GB survey file upload request")

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		return
	}

	weekNo, err := strconv.Atoi(week)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	gbInfo, err := FindGBBatch(weekNo, yearNo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	a := ws.NewFileUploads()
	si.fileUploads = a.Add(fileName)

	go func() {
		defer func() {
			si.mutux.Lock()
			si.uploadInProgress = false
			si.mutux.Unlock()
			_ = os.Remove(tmpfile)
		}()
		si.parseGBSurveyFile(tmpfile, fileName, weekNo, yearNo, gbInfo.Id)
	}()

	w.WriteHeader(http.StatusAccepted)
	OkayResponse{OK}.sendResponse(w, r)
}

func (si *SurveyImportHandler) SurveyUploadNIHandler(w http.ResponseWriter, r *http.Request) {

	si.mutux.Lock()

	w.Header().Set("Content-Type", "application/json")

	if si.uploadInProgress {
		log.Error().Msg("survey file is currently being uploaded")
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: "survey file is currently being uploaded"}.sendResponse(w, r)
		si.mutux.Unlock()
		return
	}
	si.uploadInProgress = true
	si.mutux.Unlock()

	vars := mux.Vars(r)
	month := vars["month"]
	year := vars["year"]

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received NI survey file upload request")

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		return
	}

	monthNo, err := strconv.Atoi(month)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	niInfo, err := FindNIBatch(monthNo, yearNo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	a := ws.NewFileUploads()
	si.fileUploads = a.Add(fileName)

	go func() {
		defer func() {
			si.mutux.Lock()
			si.uploadInProgress = true
			si.mutux.Unlock()
			_ = os.Remove(tmpfile)
		}()
		si.parseNISurveyFile(tmpfile, fileName, monthNo, yearNo, niInfo.Id)
	}()

	w.WriteHeader(http.StatusAccepted)
	OkayResponse{OK}.sendResponse(w, r)
}
