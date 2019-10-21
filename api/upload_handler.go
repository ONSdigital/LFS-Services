package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/util"
	"time"
)

func (h RestHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received FileUpload request")

	startTime := time.Now()

	h.w = w
	h.r = r

	vars := mux.Vars(r)
	fileType := vars["fileType"]

	w.Header().Set("Content-Type", "application/json")

	switch fileType {
	case AddressFile:
		h.uploadAddress(w, r).sendResponse(w, r)
	case SurveyFile:
		h.uploadSurvey(w, r).sendResponse(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		UnknownFileType{Status: Error, ErrorMessage: "file path in URI not recognised"}.sendResponse(w, r)
	}

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("FileUpload request completed")
}
