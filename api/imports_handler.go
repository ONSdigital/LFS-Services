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
	"strconv"
	"time"
)

const (
	Error = "ERROR"
	OK    = "OK"
)

type ImportsHandler struct{}

/*
Create a new RestHandler
*/
func NewImportsHandler() *ImportsHandler {
	return &ImportsHandler{}
}

func (im ImportsHandler) SurveyUploadGBHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	week := vars["week"]
	year := vars["year"]

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("week", week).
		Str("yeay", year).
		Msg("Received GB survey file upload request")

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		return
	}

	weekNo, err := strconv.Atoi(week)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	gbInfo, err := findGBBatch(weekNo, yearNo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	tmpfile, err := im.saveStreamToTempFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	defer func() { _ = os.Remove(tmpfile) }()

	if err := im.parseGBSurveyFile(tmpfile, fileName, weekNo, yearNo, gbInfo.Id); err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)

}

func (im ImportsHandler) SurveyUploadNIHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received NI survey file upload request")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	month := vars["month"]
	year := vars["year"]

	fileName := r.FormValue("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		ErrorResponse{Status: Error, ErrorMessage: "fileName not set"}.sendResponse(w, r)
		return
	}

	monthNo, err := strconv.Atoi(month)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	yearNo, err := strconv.Atoi(year)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	niInfo, err := findNIBatch(monthNo, yearNo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	tmpfile, err := im.saveStreamToTempFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	defer func() { _ = os.Remove(tmpfile) }()

	if err := im.parseNISurveyFile(tmpfile, fileName, monthNo, yearNo, niInfo.Id); err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}

func (im ImportsHandler) AddressUploadHandler(w http.ResponseWriter, r *http.Request) {

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

	tmpfile, err := im.saveStreamToTempFile(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	defer func() { _ = os.Remove(tmpfile) }()

	if err := im.parseAddressFile(tmpfile, fileName); err != nil {
		ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}

func (im ImportsHandler) saveStreamToTempFile(w http.ResponseWriter, r *http.Request) (string, error) {

	file, _, err := r.FormFile("lfsFile")
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

	_ = r.ParseMultipartForm(64 << 20)

	fileName := r.Form.Get("fileName")
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
