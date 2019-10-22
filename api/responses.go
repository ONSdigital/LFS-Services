package api

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Response interface {
	sendResponse(w http.ResponseWriter, r *http.Request)
}

type UnknownFileType struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
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

func (response UnknownFileType) sendResponse(w http.ResponseWriter, r *http.Request) {
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
