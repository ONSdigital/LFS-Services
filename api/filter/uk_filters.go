package filter

import (
	"fmt"
	"reflect"
	"services/types"
	"strconv"
)

type UKFilter struct {
	BaseFilter
}

func findPosition(headers []string, column string) (int, error) {
	for i, j := range headers {
		if j == column {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %s not found", column)
}

func (sf UKFilter) addHSerial(header []string, rows [][]string) (types.Column, error) {

	header = append(header, "HSERIAL")

	// get indexes of items we are interested in for the calculation
	quotaInx, err := findPosition(header, "QUOTA")
	if err != nil {
		return types.Column{}, err
	}
	weekInx, err := findPosition(header, "WEEK")
	if err != nil {
		return types.Column{}, err
	}
	w1yrInx, err := findPosition(header, "W1YR")
	if err != nil {
		return types.Column{}, err
	}
	qrtrInx, err := findPosition(header, "QRTR")
	if err != nil {
		return types.Column{}, err
	}
	addrInx, err := findPosition(header, "ADDR")
	if err != nil {
		return types.Column{}, err
	}
	wavfndInx, err := findPosition(header, "WAVFND")
	if err != nil {
		return types.Column{}, err
	}
	hhldInx, err := findPosition(header, "HHLD")
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
		Name:  "HSERIAL",
		Skip:  false,
		ColNo: len(header),
		Kind:  reflect.Int64,
	}

	return column, nil
}

func (sf UKFilter) addCaseno(header []string, rows [][]string) (types.Column, error) {

	header = append(header, "CASENO")

	// get indexes of items we are interested in for the calculation
	quotaInx, err := findPosition(header, "QUOTA")
	if err != nil {
		return types.Column{}, err
	}
	weekInx, err := findPosition(header, "WEEK")
	if err != nil {
		return types.Column{}, err
	}
	w1yrInx, err := findPosition(header, "W1YR")
	if err != nil {
		return types.Column{}, err
	}
	qrtrInx, err := findPosition(header, "QRTR")
	if err != nil {
		return types.Column{}, err
	}
	addrInx, err := findPosition(header, "ADDR")
	if err != nil {
		return types.Column{}, err
	}
	wavfndInx, err := findPosition(header, "WAVFND")
	if err != nil {
		return types.Column{}, err
	}
	hhldInx, err := findPosition(header, "HHLD")
	if err != nil {
		return types.Column{}, err
	}
	persnoInx, err := findPosition(header, "PERSNO")
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
		Name:  "CASENO",
		Skip:  false,
		ColNo: len(header),
		Kind:  reflect.Int64,
	}

	return column, nil
}
