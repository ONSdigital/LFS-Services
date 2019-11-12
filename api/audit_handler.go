package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"services/types"
	"strconv"
)

type AuditHandler struct{}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (a AuditHandler) HandleAllAuditRequest(w http.ResponseWriter, r *http.Request) {

	res, err := a.GetAllAudits()

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendDataResponse(w, r, res)
}

func (a AuditHandler) HandleYearAuditRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year := vars["year"]

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	res, err := a.GetAuditsForYear(types.Year(yearNo))

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendDataResponse(w, r, res)
}

func (a AuditHandler) HandleWeekAuditRequest(w http.ResponseWriter, r *http.Request) {

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

	res, err := a.GetAuditsForWeek(types.Week(weekNo), types.Year(yearNo))

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendDataResponse(w, r, res)
}

func (a AuditHandler) HandleMonthAuditRequest(w http.ResponseWriter, r *http.Request) {

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

	res, err := a.GetAuditsForMonth(types.Month(monthNo), types.Year(yearNo))

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendDataResponse(w, r, res)

}
