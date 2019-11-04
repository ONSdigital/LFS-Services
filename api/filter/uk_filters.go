package filter

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/types"
	"services/util"
	"strconv"
	"time"
)

type UKFilter struct {
	BaseFilter
}

func findLocation(headers []string, column string) (int, error) {
	for i, j := range headers {
		if j == column {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %s not found in findLocation()", column)
}

func (sf UKFilter) addHSerial(data [][]string) (types.Column, error) {
	header := data[0]
	rows := data[1:]

	header = append(header, "Hserial")

	// get indexes of items we are interested in for the calculation
	quotaInx, err := findLocation(header, "Quota")
	if err != nil {
		return types.Column{}, err
	}
	weekInx, err := findLocation(header, "Week")
	if err != nil {
		return types.Column{}, err
	}
	w1yrInx, err := findLocation(header, "W1Yr")
	if err != nil {
		return types.Column{}, err
	}
	qrtrInx, err := findLocation(header, "Qrtr")
	if err != nil {
		return types.Column{}, err
	}
	addrInx, err := findLocation(header, "Addr")
	if err != nil {
		return types.Column{}, err
	}
	wavfndInx, err := findLocation(header, "WavFnd")
	if err != nil {
		return types.Column{}, err
	}
	hhldInx, err := findLocation(header, "Hhld")
	if err != nil {
		return types.Column{}, err
	}

	for _, j := range rows {

		var row = j

		quota, err := strconv.ParseFloat(row[quotaInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		week, err := strconv.ParseFloat(row[weekInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		w1yr, err := strconv.ParseFloat(row[w1yrInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		qrtr, err := strconv.ParseFloat(row[qrtrInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		addr, err := strconv.ParseFloat(row[addrInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		wavfnd, err := strconv.ParseFloat(row[wavfndInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		hhld, err := strconv.ParseFloat(row[hhldInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		n := (quota * 1000000000) + (week * 10000000) + (w1yr * 1000000) +
			(qrtr * 100000) + (addr * 1000) + (wavfnd * 100) + (hhld + 1)

		row = append(row, fmt.Sprintf("%f", int64(n)))
	}

	column := types.Column{
		Name:  "Hserial",
		Skip:  false,
		ColNo: len(header),
		Kind:  reflect.Int64,
	}

	return column, nil
}

func (sf UKFilter) addCaseno(data [][]string) (types.Column, error) {

	header := data[0]
	rows := data[1:]

	startAllrows := time.Now()

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startAllrows)).
		Msg("Get all rows")

	header = append(header, "CaseNo")

	// get indexes of items we are interested in for the calculation
	quotaInx, err := findLocation(header, "Quota")
	if err != nil {
		return types.Column{}, err
	}
	weekInx, err := findLocation(header, "Week")
	if err != nil {
		return types.Column{}, err
	}
	w1yrInx, err := findLocation(header, "W1Yr")
	if err != nil {
		return types.Column{}, err
	}
	qrtrInx, err := findLocation(header, "Qrtr")
	if err != nil {
		return types.Column{}, err
	}
	addrInx, err := findLocation(header, "Addr")
	if err != nil {
		return types.Column{}, err
	}
	wavfndInx, err := findLocation(header, "WavFnd")
	if err != nil {
		return types.Column{}, err
	}
	hhldInx, err := findLocation(header, "Hhld")
	if err != nil {
		return types.Column{}, err
	}
	persnoInx, err := findLocation(header, "PersNo")
	if err != nil {
		return types.Column{}, err
	}

	for _, j := range rows {

		var row = j

		quota, err := strconv.ParseFloat(row[quotaInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		week, err := strconv.ParseFloat(row[weekInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		w1yr, err := strconv.ParseFloat(row[w1yrInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		qrtr, err := strconv.ParseFloat(row[qrtrInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		addr, err := strconv.ParseFloat(row[addrInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		wavfnd, err := strconv.ParseFloat(row[wavfndInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		hhld, err := strconv.ParseFloat(row[hhldInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		persno, err := strconv.ParseFloat(row[persnoInx], 64)
		if err != nil {
			return types.Column{}, err
		}

		n := (quota * 100000000000) + (week * 1000000000) + (w1yr * 100000000) +
			(qrtr * 10000000) + (addr * 100000) + (wavfnd * 10000) + (hhld * 100) + persno

		row = append(row, fmt.Sprintf("%f", int64(n)))
	}

	column := types.Column{
		Name:  "CaseNo",
		Skip:  false,
		ColNo: len(header),
		Kind:  reflect.Int64,
	}

	return column, nil
}
