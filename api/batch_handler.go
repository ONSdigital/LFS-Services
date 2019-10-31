package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/util"
	"strconv"
	"time"
)

type BatchHandler struct{}

func NewBatchHandler() *BatchHandler {
	return &BatchHandler{}
}

func intConversion(year string) int {
	yr, err := strconv.Atoi(year)
	if err != nil {
		return -1
	}
	return yr
}

func startLog(r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received create batch request")
}

func endLog(res error, startTime time.Time, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		log.Debug().
			Msg("Batch successfully created")
		OkayResponse{OK}.sendResponse(w, r)
	}

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Create batch complete")
}

func (b BatchHandler) CreateMonthlyBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	description := r.FormValue("description")

	// Logging
	startLog(r)

	// Convert year to int
	yr := intConversion(year)
	if yr == -1 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	// Convert month to int
	mth := intConversion(month)
	if mth == -1 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid period: %s, expected one of 1-12", month)}.sendResponse(w, r)
		return
	}

	res := b.generateMonthBatchId(mth, yr, description)
	endLog(res, startTime, w, r)
}

func (b BatchHandler) CreateQuarterlyBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]
	quarter := vars["quarter"]
	description := r.FormValue("description")

	// Logging
	startLog(r)

	// Convert year to int
	yr := intConversion(year)
	if yr == -1 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	// Strip and convert period to int
	qtr := quarter[1:]
	p, err := strconv.Atoi(qtr)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid period: %s, expected one of Q1-Q4", quarter)}.sendResponse(w, r)
		return
	}

	res := b.generateQuarterBatchId(p, yr, description)
	endLog(res, startTime, w, r)

}

func (b BatchHandler) CreateAnnualBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]
	description := r.FormValue("description")

	// Logging
	startLog(r)

	// Convert year to int
	yr := intConversion(year)
	if yr == -1 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	res := b.generateYearBatchId(yr, description)
	endLog(res, startTime, w, r)

}
