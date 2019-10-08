package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	}).Debug("Received new FileUpload Request")

	h.w = w
	h.r = r

	res := h.fileUpload()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: "ERROR", ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		OkayResponse{"OK"}.sendResponse(w, r)
	}
}

func (response OkayResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
	send(w, r, response)
}

func (response ErrorResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
	send(w, r, response)
}

func send(w http.ResponseWriter, r *http.Request, response Response) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithFields(log.Fields{
			"client": r.RemoteAddr,
			"uri":    r.RequestURI,
		}).Error("json.NewEncoder failed in FileUploadHandler")
	}
}
