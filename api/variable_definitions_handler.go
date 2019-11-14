package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"sync"
)

type VariableDefinitionsHandler struct {
	mutux            *sync.Mutex
	uploadInProgress bool // we can only handle a single upload at a time
}

func NewVariableDefinitionsHandler() *VariableDefinitionsHandler {
	return &VariableDefinitionsHandler{
		mutux:            &sync.Mutex{},
		uploadInProgress: false,
	}
}

func (vd *VariableDefinitionsHandler) setUpload(val bool) {
	vd.mutux.Lock()
	vd.uploadInProgress = val
	vd.mutux.Unlock()
}

func (vd VariableDefinitionsHandler) HandleRequestVariableUpload(w http.ResponseWriter, r *http.Request) {
	if vd.uploadInProgress {
		log.Error().Msg("Variable Definitions file is currently being uploaded")
		ErrorResponse{Status: Error, ErrorMessage: "variable definitions file is currently being uploaded"}.sendResponse(w, r)
		return
	}

	vd.setUpload(true)

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("File name not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		vd.setUpload(false)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		vd.setUpload(false)
		return
	}

	defer func() {
		vd.setUpload(false)
		_ = os.Remove(tmpfile)
	}()

	if err := vd.parseVDUpload(tmpfile, fileName); err != nil {
		log.Debug().Msg("Cannot process Variable Definitions upload")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}

func (vd VariableDefinitionsHandler) HandleRequestVariable(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	variableName := vars["variable"]

	if variableName == "" {
		log.Warn().Msg("variable not defined")
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("variable not defined")}.sendResponse(w, r)
		return
	}

	res, err := vd.getVDByVariable(variableName)

	if err != nil {
		log.Warn().Err(err).Msg("Cannot process Variable Definitions upload")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		NoRecordsFoundStatus{}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendResponse(w, r, res)

}

func (vd VariableDefinitionsHandler) HandleRequestAll(w http.ResponseWriter, r *http.Request) {

	res, err := vd.getAllVD()

	if err != nil {
		log.Error().Err(err).Msg("Get all variable definitions failed")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		NoRecordsFoundStatus{}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendResponse(w, r, res)
}
