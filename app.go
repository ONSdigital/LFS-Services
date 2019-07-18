package main

import (
	"pds-go/lfs/config"
	_ "pds-go/lfs/config"

	log "github.com/sirupsen/logrus"
	"pds-go/lfs/services"
)

func main() {

	verbose := config.Config.Database.Verbose

	log.Info("Debug Flag: ", config.Config.Debug)
	log.Debug("Database logging: ", verbose)

	res := services.GetUsers()

	for _, b := range res {
		log.Info("ID: ", b.ID)
		log.Info("User: ", b.USERNAME.String)
		log.Info("Password: ", b.PASSWORD.String)
	}

	//sub := services.GetSurveySubsample()
	//for _, b := range sub {
	//	fmt.Printf("Serial: %12.0f\n", b.SERIAL)
	//}

	//r := mux.NewRouter()
	//r.HandleFunc("/login/{user}/{password}", api.LoginHandler).Methods("GET")
	//
	//log.Fatal(http.ListenAndServe(":8000", r))
}
