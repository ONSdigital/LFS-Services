package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"services/api"
	"services/api/ws"
	"services/config"
	"services/util"
	"time"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if config.Config.LogFormat == "Terminal" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Command line flag overrides the configuration file
	debug := flag.Bool("debug", false, "sets log level to debug")

	router := mux.NewRouter()

	flag.Parse()
	if *debug || config.Config.LogLevel == "Debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		router.Use(loggingMiddleware)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().
		Str("startTime", time.Now().String()).
		Msg("LFS Imports: Starting up")

	batchHandler := api.NewBatchHandler()
	dashboardHandler := api.NewDashboardHandler()
	idHandler := api.NewIdHandler()
	surveyHandler := api.NewSurveyHandler()
	addressesHandler := api.NewAddressImportHandler()
	auditHandler := api.NewAuditHandler()
	loginHandler := api.NewLoginHandler()
	vdHandler := api.NewVariableDefinitionsHandler()
	varLabHandler := api.NewValueLabelsHandler()

	// Dashboard
	router.HandleFunc("/dashboard", dashboardHandler.HandleDashboardRequest).Methods(http.MethodGet)

	// Create New Batches Handlers
	router.HandleFunc("/batches/monthly/{year}/{month}", batchHandler.CreateMonthlyBatchHandler).Methods(http.MethodPost)
	router.HandleFunc("/batches/quarterly/{year}/{quarter}", batchHandler.CreateQuarterlyBatchHandler).Methods(http.MethodPost)
	router.HandleFunc("/batches/annual/{year}", batchHandler.CreateAnnualBatchHandler).Methods(http.MethodPost)

	// Batch info
	router.HandleFunc("/batches/display/annual/{year}", idHandler.HandleAnnualBatchIdsRequest).Methods(http.MethodGet)
	router.HandleFunc("/batches/display/quarterly/{year}/{quarter}", idHandler.HandleQuarterlyBatchIdsRequest).Methods(http.MethodGet)
	router.HandleFunc("/batches/display/monthly/{year}/{month}", idHandler.HandleMonthlyBatchIdsRequest).Methods(http.MethodGet)

	// Imports
	router.HandleFunc("/imports/survey/gb/{year}/{week}", surveyHandler.SurveyUploadGBHandler).Methods(http.MethodPost)
	router.HandleFunc("/imports/survey/ni/{year}/{month}", surveyHandler.SurveyUploadNIHandler).Methods(http.MethodPost)
	router.HandleFunc("/imports/address", addressesHandler.AddressUploadHandler).Methods(http.MethodPost)
	router.HandleFunc("/imports/variable/definitions", vdHandler.HandleRequestVariableUpload).Methods(http.MethodPost)
	// TODO: Value Labels Import -->
	router.HandleFunc("/imports/variable/definitions", varLabHandler.HandleValLabRequestlUpload).Methods(http.MethodPost)

	// Audits
	router.HandleFunc("/audits", auditHandler.HandleAllAuditRequest).Methods(http.MethodGet)
	router.HandleFunc("/audits/year/{year}", auditHandler.HandleYearAuditRequest).Methods(http.MethodGet)
	router.HandleFunc("/audits/month/{year}/{month}", auditHandler.HandleMonthAuditRequest).Methods(http.MethodGet)
	router.HandleFunc("/audits/week/{year}/{week}", auditHandler.HandleWeekAuditRequest).Methods(http.MethodGet)

	// Variable Definitions
	router.HandleFunc("/variable/definitions/{variable}", vdHandler.HandleRequestVariable).Methods(http.MethodGet)
	router.HandleFunc("/variable/definitions", vdHandler.HandleRequestAll).Methods(http.MethodGet)

	// TODO: Value Labels -->
	router.HandleFunc("/variable/definitions/{value}", varLabHandler.HandleValLabRequestValue).Methods(http.MethodGet)
	router.HandleFunc("/value/labels", varLabHandler.HandleValLabRequestAll).Methods(http.MethodGet)

	// Other
	router.HandleFunc("/login/{user}", loginHandler.LoginHandler).Methods(http.MethodGet)
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

	// we'll allow anything for now. May need or want to restrict this to just the UI when we know its endpoint
	origins := []string{"*"}
	var cors = handlers.AllowedOrigins(origins)

	handlers.CORS(cors)(router)

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
		Msg("LFS Imports: Waiting for requests")

	err = srv.ListenAndServe()
	log.Fatal().
		Err(err).
		Str("service", "LFS").
		Msgf("ListenAndServe failed")
}

func loggingMiddleware(next http.Handler) http.Handler {

	log.Info().Msg("Logging middleware registered")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().
			Str("URI:", r.RequestURI).
			Str("client", r.RemoteAddr).
			Msg("-> Received request")
		startTime := time.Now()
		next.ServeHTTP(w, r)

		log.Debug().
			Str("URI:", r.RequestURI).
			Str("elapsedTime", util.FmtDuration(startTime)).
			Msg("<- Request Completed")
	})
}
