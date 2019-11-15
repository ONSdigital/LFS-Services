package filter

import (
	"fmt"
	"services/types"
	"strconv"
)

type UKFilter struct {
	BaseFilter
}

func findPosition(data *types.SavImportData, column string) (int, error) {
	for i, j := range data.Header {
		if j.VariableName == column {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %s not found", column)
}

func (sf UKFilter) addHSerial(data *types.SavImportData) error {

	// get indexes of items we are interested in for the calculation
	quotaInx, err := findPosition(data, "QUOTA")
	if err != nil {
		return err
	}
	weekInx, err := findPosition(data, "WEEK")
	if err != nil {
		return err
	}
	w1yrInx, err := findPosition(data, "W1YR")
	if err != nil {
		return err
	}
	qrtrInx, err := findPosition(data, "QRTR")
	if err != nil {
		return err
	}
	addrInx, err := findPosition(data, "ADDR")
	if err != nil {
		return err
	}
	wavfndInx, err := findPosition(data, "WAVFND")
	if err != nil {
		return err
	}
	hhldInx, err := findPosition(data, "HHLD")
	if err != nil {
		return err
	}

	for i := 0; i < data.RowCount; i++ {
		var row = data.Rows[i]

		quota, err := strconv.ParseFloat(row.RowData[quotaInx], 64)
		if err != nil {
			return err
		}

		week, err := strconv.ParseFloat(row.RowData[weekInx], 64)
		if err != nil {
			return err
		}

		w1yr, err := strconv.ParseFloat(row.RowData[w1yrInx], 64)
		if err != nil {
			return err
		}

		qrtr, err := strconv.ParseFloat(row.RowData[qrtrInx], 64)
		if err != nil {
			return err
		}

		addr, err := strconv.ParseFloat(row.RowData[addrInx], 64)
		if err != nil {
			return err
		}

		wavfnd, err := strconv.ParseFloat(row.RowData[wavfndInx], 64)
		if err != nil {
			return err
		}

		hhld, err := strconv.ParseFloat(row.RowData[hhldInx], 64)
		if err != nil {
			return err
		}

		n := (quota * 1000000000) + (week * 10000000) + (w1yr * 1000000) +
			(qrtr * 100000) + (addr * 1000) + (wavfnd * 100) + (hhld + 1)

		row.RowData = append(row.RowData, fmt.Sprintf("%d", int64(n)))
		data.Rows[i] = row
	}

	column := types.Header{
		VariableName:        "HSERIAL",
		VariableDescription: "HSERIAL calculated column",
		VariableType:        types.TypeDouble,
		VariableLength:      8,
		VariablePrecision:   0,
		LabelName:           "",
		Drop:                false,
	}
	data.Header = append(data.Header, column)
	data.HeaderCount = data.HeaderCount + 1

	return nil
}

func (sf UKFilter) addCaseno(data *types.SavImportData) error {

	// get indexes of items we are interested in for the calculation
	quotaInx, err := findPosition(data, "QUOTA")
	if err != nil {
		return err
	}
	weekInx, err := findPosition(data, "WEEK")
	if err != nil {
		return err
	}
	w1yrInx, err := findPosition(data, "W1YR")
	if err != nil {
		return err
	}
	qrtrInx, err := findPosition(data, "QRTR")
	if err != nil {
		return err
	}
	addrInx, err := findPosition(data, "ADDR")
	if err != nil {
		return err
	}
	wavfndInx, err := findPosition(data, "WAVFND")
	if err != nil {
		return err
	}
	hhldInx, err := findPosition(data, "HHLD")
	if err != nil {
		return err
	}
	persnoInx, err := findPosition(data, "PERSNO")
	if err != nil {
		return err
	}

	for i := 0; i < data.RowCount; i++ {
		var row = data.Rows[i]

		quota, err := strconv.ParseFloat(row.RowData[quotaInx], 64)
		if err != nil {
			return err
		}

		week, err := strconv.ParseFloat(row.RowData[weekInx], 64)
		if err != nil {
			return err
		}

		w1yr, err := strconv.ParseFloat(row.RowData[w1yrInx], 64)
		if err != nil {
			return err
		}

		qrtr, err := strconv.ParseFloat(row.RowData[qrtrInx], 64)
		if err != nil {
			return err
		}

		addr, err := strconv.ParseFloat(row.RowData[addrInx], 64)
		if err != nil {
			return err
		}

		wavfnd, err := strconv.ParseFloat(row.RowData[wavfndInx], 64)
		if err != nil {
			return err
		}

		hhld, err := strconv.ParseFloat(row.RowData[hhldInx], 64)
		if err != nil {
			return err
		}

		persno, err := strconv.ParseFloat(row.RowData[persnoInx], 64)
		if err != nil {
			return err
		}

		n := (quota * 100000000000) + (week * 1000000000) + (w1yr * 100000000) +
			(qrtr * 10000000) + (addr * 100000) + (wavfnd * 10000) + (hhld * 100) + persno

		row.RowData = append(row.RowData, fmt.Sprintf("%d", int64(n)))
		data.Rows[i] = row
	}

	column := types.Header{
		VariableName:        "CASENO",
		VariableDescription: "CASENO calculated column",
		VariableType:        types.TypeDouble,
		VariableLength:      8,
		VariablePrecision:   0,
		LabelName:           "",
		Drop:                false,
	}

	data.Header = append(data.Header, column)
	data.HeaderCount = data.HeaderCount + 1

	return nil
}
