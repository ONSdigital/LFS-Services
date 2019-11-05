package filter

import (
	"github.com/rs/zerolog/log"
	"math"
	"services/types"
	"services/util"
	"strconv"
	"time"
)

type GBSurveyFilter struct {
	UKFilter
}

func NewGBSurveyFilter(audit *types.Audit) Filter {
	return GBSurveyFilter{UKFilter{BaseFilter{audit}}}
}

func (sf GBSurveyFilter) SkipRowsFilter(data [][]string) ([][]string, error) {

	header := data[0]
	rows := data[1:]

	// get indexes of items we are interested in
	sex, err := findPosition(header, "Sex")
	if err != nil {
		return nil, err
	}

	age, err := findPosition(header, "Age")
	if err != nil {
		return nil, err
	}

	indout, err := findPosition(header, "IndOut")
	if err != nil {
		return nil, err
	}

	hout, err := findPosition(header, "HOut")
	if err != nil {
		return nil, err
	}

	lstho, err := findPosition(header, "LstHO")
	if err != nil {
		return nil, err
	}

	filteredRows := make([][]string, 0, 0)
	filteredRows = append(filteredRows, header)

	for _, j := range rows {

		var row = j

		s, err := strconv.ParseFloat(row[sex], 64)
		if err != nil || math.IsNaN(s) {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because no SEX")
			continue
		}

		a, err := strconv.ParseFloat(row[age], 64)
		if err != nil || math.IsNaN(a) {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column Age is missing")
			continue
		}

		ind, err := strconv.ParseFloat(row[indout], 64)
		if err != nil || ind == 5.0 {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column INDOUT is == 5 or a Nan")
			continue
		}

		HOut, err := strconv.ParseFloat(row[hout], 64)
		if err != nil {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column HOUT is not found")
			continue
		}

		lstho, err := strconv.ParseFloat(row[lstho], 64)
		if err != nil {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column ISTHO is not found")
			continue
		}

		if HOut == 11 || HOut == 12 || HOut == 20 {
			filteredRows = append(filteredRows, j)
			continue
		}

		if (HOut == 37 && lstho == 11) ||
			(HOut == 37 && lstho == 12) ||
			(HOut == 37 && lstho == 20) ||
			(HOut == 41 && lstho == 11) ||
			(HOut == 41 && lstho == 12) ||
			(HOut == 41 && lstho == 20) {
			filteredRows = append(filteredRows, j)
			continue
		}

		// skip all other rows
		sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
		log.Debug().Msg("Dropping row because criteria not met")
	}
	return filteredRows, nil
}

func (sf GBSurveyFilter) AddVariables(data [][]string) ([]types.Column, error) {
	startTime := time.Now()

	log.Debug().
		Str("variable", "CaseNo").
		Timestamp().
		Msg("Start adding variable")

	column, err := sf.addCaseno(data)
	if err != nil {
		return nil, err
	}

	columns := make([]types.Column, 0, 0)

	log.Debug().
		Str("variable", "CaseNo").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	columns = append(columns, column)
	startTime = time.Now()

	log.Debug().
		Str("variable", "HSerial").
		Timestamp().
		Msg("Start adding variable")

	column, err = sf.addHSerial(data)
	if err != nil {
		return nil, err
	}
	columns = append(columns, column)

	log.Debug().
		Str("variable", "HSerial").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	return columns, nil
}
