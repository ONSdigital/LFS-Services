package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"sync"
)

type ValueLabelsHandler struct {
	mutux            *sync.Mutex
	uploadInProgress bool // we can only handle a single upload at a time
}

func NewValueLabelsHandler() *ValueLabelsHandler {
	return &ValueLabelsHandler{
		mutux:            &sync.Mutex{},
		uploadInProgress: false,
	}
}

func (vl *ValueLabelsHandler) setUpload(val bool) {
	vl.mutux.Lock()
	vl.uploadInProgress = val
	vl.mutux.Unlock()
}

func (vl ValueLabelsHandler) HandleValLabRequestlUpload(w http.ResponseWriter, r *http.Request) {
	if vl.uploadInProgress {
		// TODO: Value Labels: Survey file? Correct error message for Value Labels?
		log.Error().Msg("Survey file is currently being uploaded")
		ErrorResponse{Status: Error, ErrorMessage: "survey file is currently being uploaded"}.sendResponse(w, r)
		vl.setUpload(false)
		return
	}

	vl.setUpload(true)

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("File name not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		vl.setUpload(false)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		vl.setUpload(false)
		return
	}

	defer func() {
		vl.setUpload(false)
		_ = os.Remove(tmpfile)
	}()

	if err := vl.parseValLabUpload(tmpfile, fileName); err != nil {
		log.Debug().Msg("Cannot process Value Upload upload")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}

func (vl ValueLabelsHandler) HandleValLabRequestValue(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	valueName := vars["value"]

	if valueName == "" {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("value not defined")}.sendResponse(w, r)
		return
	}

	res, err := vl.getValLabByValue(valueName)

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	if res == nil {
		ErrorResponse{Status: Error, ErrorMessage: "no results returned"}.sendResponse(w, r)
	}

	SendDataResponse{}.sendDataResponse(w, r, res)

}

func (vl ValueLabelsHandler) HandleValLabRequestAll(w http.ResponseWriter, r *http.Request) {

	res, err := vl.getAllVL()

	if err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
	}

	if res == nil {
		ErrorResponse{Status: Error, ErrorMessage: "no value labels found"}.sendResponse(w, r)
	}

	SendDataResponse{}.sendDataResponse(w, r, res)

}
