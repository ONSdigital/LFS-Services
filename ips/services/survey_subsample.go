package services

import (
	"pds-go/ips/db"
	"pds-go/ips/db/model"
	"fmt"
	"log"
)

func GetSurveySubsample() []model.SURVEYSUBSAMPLE {
	var ss []model.SURVEYSUBSAMPLE
	res := db.DB.Find(&ss)

	if res.Error != nil {
		log.Fatal(fmt.Errorf("cannot read survey_subsample table %v", res))
	}

	return ss
}