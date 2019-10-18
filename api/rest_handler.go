package api

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

/*
input file names / types
*/
const (
	SurveyFile  = "Survey"
	AddressFile = "Address"
	GeogFile    = "Geog"
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
	return &RestHandlers{nil, nil}
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

	enableCors(&w)

	// TODO: Retrieve and replace static variables --->
	username := "Paul"
	password := "sucks"
	// TODO: <---
	res := h.login(username, password)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		OkayResponse{OK}.sendResponse(w, r)
	}

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		TimeDiff("elapsedTime", time.Now(), startTime).
		Msg("Login request completed")
}
