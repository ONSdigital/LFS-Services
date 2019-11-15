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

func NewGBSurveyFilter() Filter {
	return GBSurveyFilter{UKFilter{BaseFilter{}}}
}

func (sf GBSurveyFilter) SkipRowsFilter(data *types.SavImportData) error {

	// get indexes of items we are interested in
	sex, err := findPosition(data, "SEX")
	if err != nil {
		return err
	}

	age, err := findPosition(data, "AGE")
	if err != nil {
		return err
	}

	indout, err := findPosition(data, "INDOUT")
	if err != nil {
		return err
	}

	hout, err := findPosition(data, "HOUT")
	if err != nil {
		return err
	}

	lstho, err := findPosition(data, "LSTHO")
	if err != nil {
		return err
	}

	filteredRows := make([]types.Rows, 0, data.RowCount)

	for _, j := range data.Rows {

		var row = j

		s, err := strconv.ParseFloat(row.RowData[sex], 64)
		if err != nil || math.IsNaN(s) {
			log.Debug().Msg("Dropping row because no SEX")
			continue
		}

		a, err := strconv.ParseFloat(row.RowData[age], 64)
		if err != nil || math.IsNaN(a) {
			log.Debug().Msg("Dropping row because column AGE is missing")
			continue
		}

		ind, err := strconv.ParseFloat(row.RowData[indout], 64)
		if err != nil || ind == 5.0 {
			log.Debug().Msg("Dropping row because column INDOUT is == 5 or a Nan")
			continue
		}

		HOut, err := strconv.ParseFloat(row.RowData[hout], 64)
		if err != nil {
			log.Debug().Msg("Dropping row because column HOUT is not found")
			continue
		}

		lstho, err := strconv.ParseFloat(row.RowData[lstho], 64)
		if err != nil {
			log.Debug().Msg("Dropping row because column ISTHO is not found")
			continue
		}

		if HOut == 11 || HOut == 12 || HOut == 20 {
			filteredRows = append(filteredRows, types.Rows{RowData: j.RowData})
			continue
		}

		if (HOut == 37 && lstho == 11) ||
			(HOut == 37 && lstho == 12) ||
			(HOut == 37 && lstho == 20) ||
			(HOut == 41 && lstho == 11) ||
			(HOut == 41 && lstho == 12) ||
			(HOut == 41 && lstho == 20) {
			filteredRows = append(filteredRows, types.Rows{RowData: j.RowData})
			continue
		}

		// skip all other rows
		log.Debug().Msg("Dropping row because criteria not met")
	}

	data.Rows = filteredRows
	data.RowCount = len(filteredRows)

	return nil
}

func (sf GBSurveyFilter) AddVariables(data *types.SavImportData) error {
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
