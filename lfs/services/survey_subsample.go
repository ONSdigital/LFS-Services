package services

import (
	"fmt"
	"log"
	"pds-go/lfs/db"
	"pds-go/lfs/db/model"
)

func GetSurveySubsample() []model.SURVEYSUBSAMPLE {
	var ss []model.SURVEYSUBSAMPLE
	res := db.DB.Find(&ss)

	if res.Error != nil {
		log.Fatal(fmt.Errorf("cannot read survey_subsample table %v", res))
	}

	return ss
}
