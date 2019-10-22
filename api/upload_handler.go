package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		h.addressUploadHandler(w, r).sendResponse(w, r)
	case SurveyFile:
		h.surveyUploadHandler(w, r).sendResponse(w, r)
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

func (h RestHandlers) surveyUploadHandler(w http.ResponseWriter, r *http.Request) Response {
	vars := mux.Vars(r)
	runId := vars["runId"]

	if runId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn().Msg("runId not set")
		return ErrorResponse{Status: Error, ErrorMessage: "runId not set"}
	}

	tmpfile, err := h.saveStreamToTempFile()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrorResponse{Status: Error, ErrorMessage: err.Error()}
	}

	defer func() { _ = os.Remove(tmpfile) }()

	if err := h.parseInputFile(SurveyFile, tmpfile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrorResponse{Status: Error, ErrorMessage: err.Error()}
	}

	return OkayResponse{OK}
}

func (h RestHandlers) addressUploadHandler(w http.ResponseWriter, r *http.Request) Response {

	tmpfile, err := h.saveStreamToTempFile()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrorResponse{Status: Error, ErrorMessage: err.Error()}
	}

	defer func() { _ = os.Remove(tmpfile) }()

	if err := h.parseInputFile(AddressFile, tmpfile); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ErrorResponse{Status: Error, ErrorMessage: err.Error()}
	}

	return OkayResponse{OK}
}

func (h RestHandlers) saveStreamToTempFile() (string, error) {
	file, _, err := h.r.FormFile("lfsFile")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error getting formfile")
		return "", err
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	_ = h.r.ParseMultipartForm(64 << 20)

	fileName := h.r.Form.Get("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		return "", fmt.Errorf("fileName not set")
	}

	log.Debug().
		Str("fileName", fileName).
		Msg("Uploading file")

	startTime := time.Now()

	tmpfile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return "", fmt.Errorf("cannot create temporary file: %s ", err)
	}

	n, err := io.Copy(tmpfile, file)

	log.Debug().
		Str("fileName", fileName).
		Int64("bytesRead", n).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("File uploaded")

	_ = tmpfile.Close()
	return tmpfile.Name(), nil
}
