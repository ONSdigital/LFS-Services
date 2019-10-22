package api

import (
	"net/http"
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
