package filter

import (
	"github.com/rs/zerolog/log"
	"math"
	"services/types"
	"services/util"
	"time"
)

type NISurveyFilter struct {
	UKFilter
}

func NewNISurveyFilter(audit *types.Audit) types.Filter {
	return NISurveyFilter{UKFilter{BaseFilter{audit}}}
}

func (sf NISurveyFilter) SkipRow(row map[string]interface{}) bool {

	sex, ok := row["Sex"].(float64)
	if !ok || math.IsNaN(sex) {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		log.Debug().Msg("Dropping row because column Sex is missing")
		return true
	}
	age, ok := row["Age"].(float64)
	if !ok || math.IsNaN(age) {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		log.Debug().Msg("Dropping row because column Age is missing")
		return true
	}

	houtcome, ok := row["Houtcome"].(float64)
	if !ok || math.IsNaN(houtcome) {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		log.Debug().Msg("Dropping row because column Houtcome is missing")
		return true
	}

	if houtcome == 1.0 {
		row["Houtcome"] = 6.0
	}

	return false
}

func (sf NISurveyFilter) AddVariables(headers *[]string, data *[][]string) (int, error) {

	startTime := time.Now()

	log.Debug().
		Str("variable", "CaseNo").
		Timestamp().
		Msg("Start adding variables")

	if err := sf.addCaseno(headers, data); err != nil {
		return 0, err
	}

	if err := sf.addHSerial(headers, data); err != nil {
		return 0, err
	}

	log.Debug().
		Str("variable", "CaseNo").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variables")

	return 2, nil
}
