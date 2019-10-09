package api

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	SurveyFile = "Survey"
	GeogFile   = "Geog"
)

type Response interface {
	sendResponse(w http.ResponseWriter, r *http.Request)
}

type ErrorResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type OkayResponse struct {
	Status string `json:"status"`
}

type RestHandlers struct {
	log *log.Logger
	w   http.ResponseWriter
	r   *http.Request
}

func NewRestHandler(log *log.Logger) *RestHandlers {
	log.Info("Creating new RestHandler")
	return &RestHandlers{log, nil, nil}
}

func (h RestHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"client": r.RemoteAddr,
		"uri":    r.RequestURI,
	}).Debug("Received FileUpload request")

	startTime := time.Now()

	h.w = w
	h.r = r

	res := h.fileUpload()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: "ERROR", ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		a := OkayResponse{"OK"}
		sendResponse(h.w, h.r, a)
	}

	elapsed := time.Now().Sub(startTime)

	log.WithFields(log.Fields{
		"client":      r.RemoteAddr,
		"uri":         r.RequestURI,
		"elapsedTime": elapsed,
	}).Debug("FileUpload request completed")

}

func (h RestHandlers) getParameter(parameter string) (string, error) {
	keys, ok := h.r.URL.Query()[parameter]

	if !ok || len(keys[0]) < 1 {
		h.log.WithFields(log.Fields{
			"parameter": parameter,
		}).Error("URL parameter missing")
		return "", fmt.Errorf("URL parameter, %s, is missing", parameter)
	}

	return keys[0], nil
}

func (response OkayResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, response)
}

func (response ErrorResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, response)
}

func sendResponse(w http.ResponseWriter, r *http.Request, response Response) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithFields(log.Fields{
			"client": r.RemoteAddr,
			"uri":    r.RequestURI,
		}).Error("json.NewEncoder() failed in FileUploadHandler")
	}
}
