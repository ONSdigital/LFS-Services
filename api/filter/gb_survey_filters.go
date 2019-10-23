package filter

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"services/dataset"
	"services/types"
	"services/util"
	"time"
)

type GBSurveyFilter struct {
	UKFilter
}

func NewGBSurveyFilter(dataset *dataset.Dataset) types.Filter {
	return GBSurveyFilter{UKFilter{BaseFilter{dataset: dataset}}}
}

func (sf GBSurveyFilter) findLocation(headers []string, column string) (int, error) {
	for i, j := range headers {
		if j == column {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %s not found in findLoaction()", column)
}

func (sf GBSurveyFilter) SkipRow(row map[string]interface{}) bool {

	sex, ok := row["SEX"].(float64)
	if !ok || math.IsNaN(sex) {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column SEX is missing")
		return true
	}
	age, ok := row["AGE"].(float64)
	if !ok || math.IsNaN(age) {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column AGE is missing")
		return true
	}
	indout, ok := row["INDOUT"].(float64)
	if !ok || indout == 5.0 {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column INDOUT is == 5")
		return true
	}

	HOut, ok := row["HOUT"].(float64)
	if !ok {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column HOUT is not found")
		return true
	}
	lstho, ok := row["LSTHO"].(float64)
	if !ok {
		sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column ISTHO is not found")
		return true
	}

	if HOut == 11 || HOut == 12 || HOut == 20 {
		return false
	}
	if (HOut == 37 && lstho == 11) ||
		(HOut == 37 && lstho == 12) ||
		(HOut == 37 && lstho == 20) ||
		(HOut == 41 && lstho == 11) ||
		(HOut == 41 && lstho == 12) ||
		(HOut == 41 && lstho == 20) {
		return false
	}

	// skip all other rows
	sf.dataset.NumObLoaded = sf.dataset.NumObLoaded - 1
	//log.Debug().Msg("Dropping row because criteia not met")
	return true
}

func (sf GBSurveyFilter) AddVariables() (int, error) {

	startTime := time.Now()

	log.Debug().
		Str("variable", "CASENO").
		Timestamp().
		Msg("Start adding variable")

	if err := sf.addCASENO(); err != nil {
		return 0, err
	}

	log.Debug().
		Str("variable", "CASENO").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	startTime = time.Now()

	log.Debug().
		Str("variable", "HSERIAL").
		Timestamp().
		Msg("Start adding variable")

	if err := sf.addHSerial(); err != nil {
		return 0, err
	}

	log.Debug().
		Str("variable", "HSERIAL").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	return 2, nil
}
