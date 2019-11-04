package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/util"
	"time"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (d DashboardHandler) HandleDashboardRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received Dashboard request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//Functionality
	res, err := d.GetDashboardInfo()

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
			ErrorMessage: fmt.Sprintf("There is nothing to show")}.sendResponse(w, r)
		return
	}

	// Return valid json or handle
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in Dashboard")
	}

	// Logging
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retrieve Dashboard request completed")
}
