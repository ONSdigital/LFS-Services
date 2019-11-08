package api

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"services/api/ws"
	"services/types"
	"sync"
)

type AddressImportHandler struct {
	fileUploads      *types.WSMessage
	uploadInProgress bool // we can only handle a single upload to the address file at a time
	mutux            *sync.Mutex
}

func NewAddressImportHandler() *AddressImportHandler {
	return &AddressImportHandler{
		fileUploads:      nil,
		uploadInProgress: false,
		mutux:            &sync.Mutex{}}
}

func (ah *AddressImportHandler) AddressUploadHandler(w http.ResponseWriter, r *http.Request) {

	ah.mutux.Lock()

	w.Header().Set("Content-Type", "application/json")

	if ah.uploadInProgress {
		log.Error().Msg("file is currently being uploaded")
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: "address file is currently being uploaded"}.sendResponse(w, r)
		ah.mutux.Unlock()
		return
	}

	ah.uploadInProgress = true
	ah.mutux.Unlock()

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received address file upload request")

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("address upload - fileName not set")
		ErrorResponse{Status: Error, ErrorMessage: "address upload - fileName not set"}.sendResponse(w, r)
		return
	}

	tmpfile, err := SaveStreamToTempFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	a := ws.NewFileUploads()
	ah.fileUploads = a.Add(fileName)

	go func() {
		defer func() {
			ah.mutux.Lock()
			ah.uploadInProgress = false
			ah.mutux.Unlock()
			_ = os.Remove(tmpfile)
		}()
		ah.ParseAddressFile(tmpfile, fileName)
	}()

	w.WriteHeader(http.StatusAccepted)
	OkayResponse{OK}.sendResponse(w, r)
}
