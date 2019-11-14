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

func (sf GBSurveyFilter) SkipRowsFilter(header []string, rows [][]string, data types.SavImportData) ([][]string, error) {

	// get indexes of items we are interested in
	sex, err := findPosition(header, "SEX")
	if err != nil {
		return nil, err
	}

	age, err := findPosition(header, "AGE")
	if err != nil {
		return nil, err
	}

	indout, err := findPosition(header, "INDOUT")
	if err != nil {
		return nil, err
	}

	hout, err := findPosition(header, "HOUT")
	if err != nil {
		return nil, err
	}

	lstho, err := findPosition(header, "LSTHO")
	if err != nil {
		return nil, err
	}

	filteredRows := make([][]string, 0, 0)

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
			log.Debug().Msg("Dropping row because column AGE is missing")
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

func (sf GBSurveyFilter) AddVariables(headers []string, rows [][]string, data types.SavImportData) ([]types.Column, error) {
	startTime := time.Now()

	log.Debug().
		Str("variable", "CASENO").
		Timestamp().
		Msg("Start adding variable")

	column, err := sf.addCaseno(headers, rows, data)
	if err != nil {
		return nil, err
	}

	columns := make([]types.Column, 0, 0)

	log.Debug().
		Str("variable", "CASENO").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	columns = append(columns, column)
	startTime = time.Now()

	log.Debug().
		Str("variable", "HSERIAL").
		Timestamp().
		Msg("Start adding variable")

	column, err = sf.addHSerial(headers, rows, data)
	if err != nil {
		return nil, err
	}
	columns = append(columns, column)

	log.Debug().
		Str("variable", "HSERIAL").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	return columns, nil
}
