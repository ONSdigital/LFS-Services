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
		Msg("Received get Annual Batch ID request")

	// Convert year to integer
	yearNo, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Functionality
	res, err := i.GetIdsForYear(types.Year(yearNo))

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
	sendIdResponse(w, r, res)

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retrieve Annual Batch ID request completed")
}

func sendIdResponse(w http.ResponseWriter, r *http.Request, response []types.YearID) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in sendIdResponse")
	}
}
