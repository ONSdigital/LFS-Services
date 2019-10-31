package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"services/util"
	"time"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (l LoginHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("client", r.RemoteAddr).
		Str("uri", r.RequestURI).
		Msg("Received login request")

	startTime := time.Now()

	// Assign username and password variables
	vars := mux.Vars(r)
	username := vars["user"]
	password := r.Header.Get("password")

	// Call login service to validate
	res := l.login(username, password)

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
