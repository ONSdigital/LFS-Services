package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"services/api"
	"services/api/ws"
	"services/config"
	"time"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if config.Config.LogFormat == "Terminal" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if config.Config.LogLevel == "Debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Command line flag overrides the configuration file
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().
		Str("startTime", time.Now().String()).
		Msg("LFS Services: Starting up")

	router := mux.NewRouter()
	restHandlers := api.NewRestHandler()

	router.HandleFunc("/import/{fileType}/{runId}", restHandlers.FileUploadHandler).Methods(http.MethodPost)
	router.HandleFunc("/import/{fileType}", restHandlers.FileUploadHandler).Methods(http.MethodPost)
	router.HandleFunc("/login/{user}", restHandlers.LoginHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", ws.WebSocketHandler{}.ServeWs).Methods(http.MethodGet)

	listenAddress := config.Config.Service.ListenAddress

	writeTimeout, err := time.ParseDuration(config.Config.Service.WriteTimeout)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "LFS").
			Msgf("writeTimeout configuration error")
	}

	readTimeout, err := time.ParseDuration(config.Config.Service.ReadTimeout)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "LFS").
			Msgf("readTimeout configuration error")
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         listenAddress,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	log.Info().
		Str("listenAddress", listenAddress).
		Str("writeTimeout", writeTimeout.String()).
		Str("readTimeout", readTimeout.String()).
		Msg("LFS Services: Waiting for requests")

	err = srv.ListenAndServe()
	log.Fatal().
		Err(err).
		Str("service", "LFS").
		Msgf("ListenAndServe failed")
}
