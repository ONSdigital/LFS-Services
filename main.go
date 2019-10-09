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
	logger.Info("LFS Services: Starting up")

	router.HandleFunc("/import/survey/{run_id}", restHandlers.FileUploadHandler).Methods("POST")

	listenAddress := config.Config.Service.ListenAddress

	writeTimeout, err := time.ParseDuration(config.Config.Service.WriteTimeout)
	if err != nil {
		panic("writeTimeout configuration error")
	}

	readTimeout, err := time.ParseDuration(config.Config.Service.ReadTimeout)
	if err != nil {
		panic("readTimeout configuration error")
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         listenAddress,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	logger.WithFields(log.Fields{
		"listenAddress": listenAddress,
		"writeTimeout":  writeTimeout,
		"readTimeout":   readTimeout,
	}).Info("LFS Services: Waiting for requests")

	logger.Fatal(srv.ListenAndServe())
}
