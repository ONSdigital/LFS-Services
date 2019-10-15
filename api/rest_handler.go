package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

/*
input file names / types
*/
const (
	SurveyFile = "Survey"
	GeogFile   = "Geog"
)

const (
	Error = "ERROR"
	OK    = "OK"
)

type RestHandlers struct {
	w http.ResponseWriter
	r *http.Request
}

/*
Create a new RestHandler
*/
func NewRestHandler() *RestHandlers {
	log.Debug().
		Time("startTime", time.Now()).
		Msg("Creating new RestHandler")
	return &RestHandlers{nil, nil}
}

func (h RestHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received FileUpload request")

	startTime := time.Now()

	h.w = w
	h.r = r

	res := h.fileUpload()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		a := OkayResponse{OK}
		sendResponse(h.w, h.r, a)
	}

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("FileUpload request completed")

}

func (h RestHandlers) getParameter(parameter string) (string, error) {
	keys, ok := h.r.URL.Query()[parameter]

	if !ok || len(keys[0]) < 1 {
		log.Error().
			Str("parameter", parameter).
			Msg("URL parameter missing")
		return "", fmt.Errorf("URL parameter, %s, is missing", parameter)
	}

	return keys[0], nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (h RestHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received login request")

	startTime := time.Now()

	h.w = w
	h.r = r

	// TODO: Retrieve and replace static variables --->
	username := "Paul"
	password := "sucks"
	// TODO: <---
	res := h.login(username, password)

	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		a := OkayResponse{OK}
		sendResponse(h.w, h.r, a)
	}

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("Login request completed")

}

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

func (response OkayResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, response)
}

func (response ErrorResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, response)
}

func sendResponse(w http.ResponseWriter, r *http.Request, response Response) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().
			Str("client", r.RemoteAddr).
			Str("uri", r.RequestURI).
			Msg("json.NewEncoder() failed in FileUploadHandler")
	}
}
