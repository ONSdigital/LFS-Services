package filter

import (
	"github.com/rs/zerolog/log"
	"math"
	"services/dataset"
	"services/types"
	"services/util"
	"time"
)

type NISurveyFilter struct {
	UKFilter
}

func NewNISurveyFilter(dataset *dataset.Dataset) types.Filter {
	return NISurveyFilter{UKFilter{BaseFilter{dataset: dataset}}}
}

func (sf NISurveyFilter) SkipRow(row map[string]interface{}) bool {

	sex, ok := row["SEX"].(float64)
	if !ok || math.IsNaN(sex) {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		log.Debug().Msg("Dropping row because column SEX is missing")
		return true
	}
	age, ok := row["AGE"].(float64)
	if !ok || math.IsNaN(age) {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		log.Debug().Msg("Dropping row because column AGE is missing")
		return true
	}

	houtcome, ok := row["HOUTCOME"].(float64)
	if !ok || math.IsNaN(houtcome) {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		log.Debug().Msg("Dropping row because column HOUTCOME is missing")
		return true
	}

	if houtcome == 1.0 {
		row["HOUTCOME"] = 6.0

	}

	return false
}

func (sf NISurveyFilter) AddVariables() (int, error) {

	startTime := time.Now()

	log.Debug().
		Str("variable", "CASENO").
		Timestamp().
		Msg("Start adding variables")

	if err := sf.addCASENO(); err != nil {
		return 0, err
	}

	if err := sf.addHSerial(); err != nil {
		return 0, err
	}

	log.Debug().
		Str("variable", "CASENO").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variables")

	return 2, nil
}
