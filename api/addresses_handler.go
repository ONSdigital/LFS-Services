package api

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"services/api/ws"
	"services/types"
)

type AddressImportHandler struct {
	fileUploads *types.WSMessage
}

func NewAddressImportHandler() *AddressImportHandler {
	return &AddressImportHandler{nil}
}

func (ah AddressImportHandler) AddressUploadHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received address file upload request")

	w.Header().Set("Content-Type", "application/json")

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

	defer func() { _ = os.Remove(tmpfile) }()

	a := ws.NewFileUploads()
	ah.fileUploads = a.Add(fileName)

	go func() {
		ah.ParseAddressFile(tmpfile, fileName)
	}()

	OkayResponse{OK}.sendResponse(w, r)
}
