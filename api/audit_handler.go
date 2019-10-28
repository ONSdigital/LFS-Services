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

type AuditHandler struct{}

/*
Create a new RestHandler
*/
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (h AuditHandler) HandleAllAuditRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received get all audit request")

	startTime := time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := h.GetAllAudits()

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	sendAuditResponse(w, r, res)

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retieve audit request completed")
}

func (h AuditHandler) HandleYearAuditRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received get year audit request")

	startTime := time.Now()

	vars := mux.Vars(r)
	year := vars["year"]

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := h.GetAuditsForYear(types.Year(yearNo))

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	sendAuditResponse(w, r, res)

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retieve audit request completed")
}

func (h AuditHandler) HandleWeekAuditRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received get week audit request")

	startTime := time.Now()

	vars := mux.Vars(r)
	year := vars["year"]
	week := vars["week"]

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	weekNo, err := strconv.Atoi(week)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid week: %s, expected an integer", week)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := h.GetAuditsForWeek(types.Week(weekNo), types.Year(yearNo))

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	sendAuditResponse(w, r, res)

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retieve audit request completed")
}

func (h AuditHandler) HandleMonthAuditRequest(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received get month audit request")

	startTime := time.Now()

	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	monthNo, err := strconv.Atoi(month)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid month: %s, expected an integer", month)}.sendResponse(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := h.GetAuditsForMonth(types.Month(monthNo), types.Year(yearNo))

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	sendAuditResponse(w, r, res)

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Retieve audit request completed")
}

func sendAuditResponse(w http.ResponseWriter, r *http.Request, response []types.Audit) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in sendAuditResponse")
	}
}
