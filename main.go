package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"services/api"
	"services/config"
	"time"
)

func main() {
	router := mux.NewRouter()
	logger := log.New()
	restHandlers := api.NewRestHandler(logger)
	logger.Info("Starting up")
	router.HandleFunc("/import/survey/{run_id}", restHandlers.FileUploadHandler).Methods("POST")

	listenAddress := config.Config.Service.ListenAddress

	srv := &http.Server{
		Handler:      router,
		Addr:         listenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Waiting for requests")
	logger.Fatal(srv.ListenAndServe())
}
