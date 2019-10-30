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

func (b BatchHandler) CreateBatchHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received create batch request")

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
	description := r.FormValue("description")

	switch period {
	case "0":
		res = b.handleYear(yr, description)
	case "Q1", "Q2", "Q3", "Q4":
		res = b.handleQuarter(period, yr)
	default:
		{
			month, err := strconv.Atoi(period)
			if err == nil {
				res = b.handleMonth(month, yr, description)
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
		Msg("Create batch complete")
}

func (b BatchHandler) handleMonth(month int, year int, description string) error {
	res := b.generateMonthBatchId(month, year, description)
	return res
}

func (b BatchHandler) handleQuarter(quarter string, year int) error {
	//res := b.generateQuarterBatchId(quarter, year)
	return nil
}

func (b BatchHandler) handleYear(year int, description string) error {
	res := b.generateYearBatchId(year, description)
	return res
}
