package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	SURVEY_FILE = "Survey"
	GEOG_FILE   = "Geog"
)

type Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"status"`
	Request      string `json:"request"`
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

func (handler RestHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"client": r.RemoteAddr,
		"uri":    r.RequestURI,
	}).Debug("Received new FileUpload Request")

	handler.w = w
	handler.r = r

	res := handler.fileUpload()

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	var stat Response

	if res == nil {
		stat = Response{
			"OK",
			res.Error(),
			vars["run_id"],
		}
	} else {
		stat = Response{
			"OK",
			nil,
			vars["run_id"],
		}
	}

	if err := json.NewEncoder(w).Encode(stat); err != nil {
		log.WithFields(log.Fields{
			"client": r.RemoteAddr,
			"uri":    r.RequestURI,
		}).Error("json.NewEncoder failed in FileUploadHandler")
	}
}
