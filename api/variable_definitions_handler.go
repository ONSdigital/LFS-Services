package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type VariableDefinitionsHandler struct{}

func NewVariableDefinitionsHandler() *VariableDefinitionsHandler {
	return &VariableDefinitionsHandler{}
}

func (vd VariableDefinitionsHandler) HandleRequestVariable(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	variableName := vars["variable"]

	if variableName == "" {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("variable not defined")}.sendResponse(w, r)
		return
	}

	res, err := vd.getVDByVariable(variableName)

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	if res == nil {
		ErrorResponse{Status: Error, ErrorMessage: "no results returned"}.sendResponse(w, r)
	}

	SendDataResponse{}.sendDataResponse(w, r, res)

}

func (vd VariableDefinitionsHandler) HandleRequestAll(w http.ResponseWriter, r *http.Request) {

	res, err := vd.getAllVD()

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	if res == nil {
		ErrorResponse{Status: Error, ErrorMessage: "no variable definitions found"}.sendResponse(w, r)
	}

	SendDataResponse{}.sendDataResponse(w, r, res)

}
