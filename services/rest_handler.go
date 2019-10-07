package services

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Request string `json:"request"`
}

type RestHandlers struct {
	log *log.Logger
}

func NewRestHandler(log *log.Logger) *RestHandlers {
	log.Info("Creating new RestHandler")
	return &RestHandlers{log}
}

func (handler RestHandlers) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"client": r.RemoteAddr,
		"uri":    r.RequestURI,
	}).Debug("Received new FileUpload Request")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	stat := Response{"OK", vars["run_id"]}

	if err := json.NewEncoder(w).Encode(stat); err != nil {
		log.WithFields(log.Fields{
			"client": r.RemoteAddr,
			"uri":    r.RequestURI,
		}).Error("json.NewEncoder failed in FileUploadHandler")
	}
}
