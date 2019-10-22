package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/util"
	"time"
)

/*
input file names / types
*/
const (
	SurveyFile  = "Survey"
	AddressFile = "Address"
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

	// Assign username and password variables
	vars := mux.Vars(r)
	username := vars["user"]
	password := h.r.Header.Get("password")

	// Call login function to validate
	res := h.login(username, password)

	// Enable "Cross-Origin Resource Sharing"
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if res != nil {
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
	} else {
		log.Debug().
			Msg("Login request successful")
		OkayResponse{OK}.sendResponse(w, r)
	}

	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Login request completed")
}
