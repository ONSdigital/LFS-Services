package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BatchHandler struct {
}

func NewBatchHandler() *BatchHandler {
	return &BatchHandler{}
}

func (b BatchHandler) CreateMonthlyBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables

	vars := mux.Vars(r)
	year := vars["year"]
	month := vars["month"]
	description := r.FormValue("description")

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
	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)

}

func (b BatchHandler) CreateQuarterlyBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables
	vars := mux.Vars(r)
	year := vars["year"]
	quarter := vars["quarter"]
	description := r.FormValue("description")

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
	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}

func (b BatchHandler) CreateAnnualBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables
	vars := mux.Vars(r)
	year := vars["year"]
	description := r.FormValue("description")

	// Convert year to int
	yr := intConversion(year)
	if yr == -1 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	res := b.generateYearBatchId(yr, description)
	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}
