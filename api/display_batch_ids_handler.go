package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/types"
	"services/util"
	"strconv"
	"time"
)

type IdHandler struct{}

func NewIdHandler() *IdHandler {
	return &IdHandler{}
}

func (i IdHandler) HandleAnnualBatchIdsRequest(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received Annual Batch ID request")

	// Convert year to integer
	yr, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Functionality
	res, err := i.GetIdsForYear(types.Year(yr))

	// Error handling
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("No valid annual batches for %s", year)}.sendResponse(w, r)
		return
	}

	// Return valid json or handle
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in sendIdResponse")
	}

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retrieve Annual Batch ID request completed")
}

func (i IdHandler) HandleQuarterlyBatchIdsRequest(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]
	quarter := vars["quarter"]

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received Quarterly Batch ID request")

	// Convert year to integer
	yr, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, %s, expected an integer", quarter, year)}.sendResponse(w, r)
		return
	}

	// Strip and convert quarter to int
	qtr := quarter[1:]
	q, err := strconv.Atoi(qtr)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid period: %s, expected one of Q1-Q4", quarter)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Functionality
	res, err := i.GetIdsForQuarter(types.Year(yr), types.Quarter(q))

	// Error handling
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("No valid quarter batches for %s", year)}.sendResponse(w, r)
		return
	}

	// Return valid json or handle
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in sendIdResponse")
	}

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retrieve Quarterly Batch ID request completed")
}

func (i IdHandler) HandleMonthlyBatchIdsRequest(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received Monthly Batch ID request")

	// Convert year to integer
	yr, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	// Convert month to integer
	mth, err := strconv.Atoi(month)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid month: %s, expected one of 1-12", month)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Functionality
	res, err := i.GetIdsForMonth(types.Year(yr), types.Month(mth))

	// Error handling
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("No valid monthly batches for %s", year)}.sendResponse(w, r)
		return
	}

	// Return valid json or handle
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in sendIdResponse")
	}

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retrieve Monthly Batch ID request completed")
}

func (i IdHandler) HandleNIBatchIdsRequest(w http.ResponseWriter, r *http.Request) {
	// Variables
	startTime := time.Now()
	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received NI Batch ID request")

	// Convert year to integer
	yr, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	// Convert month to integer
	mth, err := strconv.Atoi(month)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid month: %s, expected one of 1-12", month)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Functionality
	res, err := i.GetIdsForNI(types.Year(yr), types.Month(mth))

	// Error handling
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("No valid NI batches for %s", year)}.sendResponse(w, r)
		return
	}

	// Return valid json or handle
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in sendIdResponse")
	}

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retrieve NI Batch IDs request completed")
}
