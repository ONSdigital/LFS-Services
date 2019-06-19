package main

import (
	"pds-go/ips/config"
	_ "pds-go/ips/config"

	"pds-go/ips/services"
	log "github.com/sirupsen/logrus"
)

func main() {

	verbose := config.Config.Database.Verbose

	log.Info("Debug Flag: ", config.Debug)
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
