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
	q, err := strconv.Atoi(quarter[1:])
	if err != nil || len(quarter) != 2 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid period: %s, expected one of Q1-Q4", quarter)}.sendResponse(w, r)
		return
	}

	// Do
	res, qErr := b.generateQuarterBatchId(q, yr, description)
	if res != nil {
		BadDataResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("3 valid months for Q%d, %d required", q, yr),
			Result:       res,
		}.sendResponse(w, r)
		return
	}
	if qErr != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: qErr.Error()}.sendResponse(w, r)
		return
	}

	// Return
	OkayResponse{OK}.sendResponse(w, r)
}

func (b BatchHandler) CreateAnnualBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Variables
	vars := mux.Vars(r)
	year := vars["year"]
	description := r.FormValue("description")

	// Convert year to int
	yr := intConversion(year)
	if yr < -1 || yr == 0 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
		return
	}

	// Do
	res, aErr := b.generateYearBatchId(yr, description)
	if res != nil {
		BadDataResponse{}.sendResponse(w, r)
		return
	}
	if aErr != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: aErr.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}
