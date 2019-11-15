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

func NewNISurveyFilter() Filter {
	return NISurveyFilter{UKFilter{BaseFilter{}}}
}

func (sf NISurveyFilter) SkipRowsFilter(data *types.SavImportData) error {

	// get indexes of items we are interested in
	sex, err := findPosition(data, "SEX")
	if err != nil {
		return err
	}

	age, err := findPosition(data, "AGE")
	if err != nil {
		return err
	}
	houtcome, err := findPosition(data, "HOUTCOME")
	if err != nil {
		return err
	}

	filteredRows := make([]types.Rows, 0, data.RowCount)

	for _, j := range data.Rows {

		var row = j

		s, err := strconv.ParseFloat(row.RowData[sex], 64)
		if err != nil {
			log.Error().Msg("SEX field is not a float, ignoring")
			continue
		}
		if math.IsNaN(s) {
			log.Debug().Msg("Dropping row because column SEX is missing")
			continue
		}

		a, err := strconv.ParseFloat(row.RowData[age], 64)
		if err != nil {
			log.Error().Msg("AGE field is not a float, ignoring")
			continue
		}
		if math.IsNaN(a) {
			log.Debug().Msg("Dropping row because column Age is missing")
			continue
		}

		h, err := strconv.ParseFloat(row.RowData[houtcome], 64)
		if err != nil {
			log.Error().Msg("HOUTCOME field is not a float, ignoring")
			continue
		}
		if math.IsNaN(h) {
			log.Debug().Msg("Dropping row because column Houtcome is missing")
			continue
		}
		filteredRows = append(filteredRows, types.Rows{RowData: j.RowData})
	}
	data.Rows = filteredRows
	data.RowCount = len(filteredRows)
	return nil
}

func (sf NISurveyFilter) AddVariables(data *types.SavImportData) error {
	startTime := time.Now()

	log.Debug().
		Str("variable", "CASENO").
		Timestamp().
		Msg("Start adding variable")

	if err := sf.addCaseno(data); err != nil {
		return err
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

	if err := sf.addHSerial(data); err != nil {
		return err
	}

	log.Debug().
		Str("variable", "HSERIAL").
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Finished adding variable")

	return nil
}
