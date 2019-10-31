package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
)

type SurveyImportHandler struct{}

func NewSurveyHandler() *SurveyImportHandler {
	return &SurveyImportHandler{}
}

func (si SurveyImportHandler) SurveyUploadGBHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

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

	defer func() { _ = os.Remove(tmpfile) }()

	if err := si.parseGBSurveyFile(tmpfile, fileName, weekNo, yearNo, gbInfo.Id); err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)

}

func (si SurveyImportHandler) SurveyUploadNIHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received NI survey file upload request")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	month := vars["month"]
	year := vars["year"]

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

	defer func() { _ = os.Remove(tmpfile) }()

	if err := si.parseNISurveyFile(tmpfile, fileName, monthNo, yearNo, niInfo.Id); err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}
