package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"services/types"
	"strings"
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
		log.Error().Msg("Value labels file is currently being uploaded")
		ErrorResponse{Status: Error, ErrorMessage: "value labels file is currently being uploaded"}.sendResponse(w, r)
		return
	}

	vars := mux.Vars(r)
	source := vars["source"]
	source = strings.ToUpper(source)

	if source == "" || (source != string(types.GBSource) && source != string(types.NISource)) {
		log.Warn().Msg("source must be either gb or ni")
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("source not defined")}.sendResponse(w, r)
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

	if err := vl.parseValLabUpload(tmpfile, fileName, types.FileSource(source)); err != nil {
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
		log.Error().Err(err).
			Str("variable", valueName).
			Msg("Get value label failed")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		log.Debug().Msg("No value labels found")
		NoRecordsFoundStatus{}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendResponse(w, r, res)

}

func (vl ValueLabelsHandler) HandleValLabRequestAll(w http.ResponseWriter, r *http.Request) {

	res, err := vl.getAllVL()

	if err != nil {
		log.Error().Err(err).Msg("Get all value labels failed")
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		log.Debug().Msg("No value labels found")
		NoRecordsFoundStatus{}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendResponse(w, r, res)

}
