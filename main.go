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

	log.WithFields(log.Fields{
		"startTime": time.Now(),
	}).Info("LFS Services: Starting up")

	router := mux.NewRouter()
	restHandlers := api.NewRestHandler()

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

	log.WithFields(log.Fields{
		"listenAddress": listenAddress,
		"writeTimeout":  writeTimeout,
		"readTimeout":   readTimeout,
	}).Info("LFS Services: Waiting for requests")

	log.Fatal(srv.ListenAndServe())
}
