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

//func (b BatchHandler) CreateBatchHandler(w http.ResponseWriter, r *http.Request) {
//	log.Debug().
//		Str("client", r.RemoteAddr).
//		Str("uri", r.RequestURI).
//		Msg("Received create batch request")
//
//	startTime := time.Now()
//
//	vars := mux.Vars(r)
//	year := vars["year"]
//	period := vars["period"]
//	description := r.FormValue("description")
//
//	// Convert year to int
//	yr, err := strconv.Atoi(year)
//	if err != nil {
//		ErrorResponse{
//			Status:       Error,
//			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
//		return
//	}
//
//	var res error
//
//	switch period {
//	case "0":
//		res = b.handleYear(yr, description)
//	case "Q1", "Q2", "Q3", "Q4":
//		// Strip and convert period to int
//		period = period[1:]
//		p, err := strconv.Atoi(period)
//		if err != nil {
//			ErrorResponse{
//				Status:       Error,
//				ErrorMessage: fmt.Sprintf("invalid period: %s, expected on of Q1-Q4", period)}.sendResponse(w, r)
//			return
//		}
//
//		res = b.handleQuarter(p, yr, description)
//	default:
//		{
//			month, err := strconv.Atoi(period)
//			if err == nil {
//				res = b.handleMonth(month, yr, description)
//				break
//			}
//			ErrorResponse{
//				Status:       Error,
//				ErrorMessage: fmt.Sprintf("invalid period: %s, expected one of 1-12, Q1-Q4", period)}.sendResponse(w, r)
//		}
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//
//	if res != nil {
//		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
//	} else {
//		log.Debug().
//			Msg("Batch successfully created")
//		OkayResponse{OK}.sendResponse(w, r)
//	}
//
//	log.Debug().
//		Str("client", r.RemoteAddr).
//		Str("uri", r.RequestURI).
//		Str("elapsedTime", util.FmtDuration(startTime)).
//		Msg("Create batch complete")
//}

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

//func (b BatchHandler) handleMonth(month int, year int, description string) error {
//	res := b.generateMonthBatchId(month, year, description)
//	return res
//}

func (b BatchHandler) handleQuarter(quarter int, year int, description string) error {
	res := b.generateQuarterBatchId(quarter, year, description)
	return res
}

func (b BatchHandler) handleYear(year int, description string) error {
	res := b.generateYearBatchId(year, description)
	return res
}
