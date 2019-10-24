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

func (h RestHandlers) CreateBatchHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received new Batch ID request")

	startTime := time.Now()

	vars := mux.Vars(r)
	year := vars["year"]
	period := vars["period"]

	yr, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	var res error

	switch period {
	case "":
		res = h.handleYear(yr)
	case "Q1", "Q2", "Q3", "Q4":
		res = h.handleQuarter(period, yr)
	default:
		{
			month, err := strconv.Atoi(period)
			description := r.FormValue("description")
			if err == nil {
				res = h.handleMonth(month, yr, description)
				break
			}
			ErrorResponse{
				Status:       Error,
				ErrorMessage: fmt.Sprintf("invalid period: %s, expected one of 1-12, Q1-Q4", period)}.sendResponse(w, r)
		}
	}

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
		Msg("Batch ID request complete")
}

func (h RestHandlers) handleMonth(month int, year int, description string) error {
	res := h.generateMonthBatchId(month, year, description)
	return res
}

func (h RestHandlers) handleQuarter(quarter string, year int) error {
	//res := h.generateQuarterBatchId(quarter, year)
	return nil
}

func (h RestHandlers) handleYear(year int) error {
	//res := h.generateYearBatchId(year)
	return nil
}
