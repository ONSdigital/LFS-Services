package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (l LoginHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Assign username and password variables
	vars := mux.Vars(r)
	username := vars["user"]
	password := r.Header.Get("password")

	// Call login service to validate
	res := l.login(username, password)

	if res != nil {
		log.Debug().Msg("Login request failed")
		ErrorResponse{Status: Error, ErrorMessage: res.Error()}.sendResponse(w, r)
		return
	}

	OkayResponse{OK}.sendResponse(w, r)
}
