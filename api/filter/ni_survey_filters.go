package filter

import (
	"github.com/rs/zerolog/log"
	"math"
	"services/types"
	"services/util"
	"strconv"
	"time"
)

type NISurveyFilter struct {
	UKFilter
}

func NewNISurveyFilter(audit *types.Audit) Filter {
	return NISurveyFilter{UKFilter{BaseFilter{audit}}}
}

func (sf NISurveyFilter) SkipRowsFilter(header []string, data [][]string) ([][]string, error) {

	// get indexes of items we are interested in
	sex, err := findPosition(header, "Sex")
	if err != nil {
		return nil, err
	}

	age, err := findPosition(header, "Age")
	if err != nil {
		return nil, err
	}
	houtcome, err := findPosition(header, "Houtcome")
	if err != nil {
		return nil, err
	}

	filteredRows := make([][]string, 0, 0)
	filteredRows = append(filteredRows, header)

	for _, j := range data {

		var row = j

		s, err := strconv.ParseFloat(row[sex], 64)
		if err != nil {
			log.Error().Msg("sex field is not a float, ignoring")
			continue
		}
		if math.IsNaN(s) {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column Sex is missing")
			continue
		}

		a, err := strconv.ParseFloat(row[age], 64)
		if err != nil {
			log.Error().Msg("age field is not a float, ignoring")
			continue
		}
		if math.IsNaN(a) {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column Age is missing")
			continue
		}

		h, err := strconv.ParseFloat(row[houtcome], 64)
		if err != nil {
			log.Error().Msg("age field is not a float, ignoring")
			continue
		}
		if math.IsNaN(h) {
			sf.Audit.NumObLoaded = sf.Audit.NumObLoaded - 1
			log.Debug().Msg("Dropping row because column Houtcome is missing")
			continue
		}
		filteredRows = append(filteredRows, j)
	}
	return filteredRows, nil
}

func (sf NISurveyFilter) AddVariables(headers []string, data [][]string) ([]types.Column, error) {
	startTime := time.Now()

	log.Debug().
		Str("variable", "CaseNo").
		Timestamp().
		Msg("Start adding variable")

	column, err := sf.addCaseno(headers, data)
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

	column, err = sf.addHSerial(headers, data)
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
