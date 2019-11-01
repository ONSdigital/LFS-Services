package filter

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"services/types"
	"services/util"
	"time"
)

type GBSurveyFilter struct {
	UKFilter
}

func NewGBSurveyFilter(audit *types.Audit) types.Filter {
	return GBSurveyFilter{UKFilter{BaseFilter{audit}}}
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
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column SEX is missing")
		return true
	}
	age, ok := row["AGE"].(float64)
	if !ok || math.IsNaN(age) {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column AGE is missing")
		return true
	}
	indout, ok := row["INDOUT"].(float64)
	if !ok || indout == 5.0 {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column INDOUT is == 5")
		return true
	}

	HOut, ok := row["HOUT"].(float64)
	if !ok {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		//log.Debug().Msg("Dropping row because column HOUT is not found")
		return true
	}
	lstho, ok := row["LSTHO"].(float64)
	if !ok {
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
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
	sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
	//log.Debug().Msg("Dropping row because criteia not met")
	return true
}

func (sf GBSurveyFilter) AddVariables(headers *[]string, data *[][]string) (int, error) {

	startTime := time.Now()

	log.Debug().
		Str("variable", "CaseNo").
		Timestamp().
		Msg("Start adding variable")

	if err := sf.addCASENO(headers, data); err != nil {
		return 0, err
	}

	log.Debug().
		Str("variable", "CaseNo").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	startTime = time.Now()

	log.Debug().
		Str("variable", "HSerial").
		Timestamp().
		Msg("Start adding variable")

	if err := sf.addHSerial(headers, data); err != nil {
		return 0, err
	}

	log.Debug().
		Str("variable", "HSerial").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	return 2, nil
}
