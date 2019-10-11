package filter

import (
	"fmt"
	"reflect"
	conf "services/config"
	"services/dataset"
	"strconv"
)

func NewSurveyFilter(dataset *dataset.Dataset) Filter {
	return SurveyFilter{dataset: dataset}
}

type SurveyFilter struct {
	dataset *dataset.Dataset
}

func (sf SurveyFilter) DropColumns() {
	drop := conf.Config.DropColumns.Survey
	_ = sf.dataset.DropColumns(drop.ColumnNames)
}

func (sf SurveyFilter) RenameColumns() {
	cols := conf.Config.Rename.Survey
	m := make(map[string]string, sf.dataset.ColumnCount)

	for _, v := range cols {
		m[v.From] = v.To
	}
	_ = sf.dataset.RenameColumns(m)
}

func findLocation(headers []string, column string) (int, error) {
	for i, j := range headers {
		if j == column {
			return i, nil
		}
	}
	return 0, fmt.Errorf("column %s not found in findLoaction()", column)
}

func (sf SurveyFilter) AddVariables() error {

	column, err := sf.dataset.AddColumn("CASENO", reflect.Int64)
	if err != nil {
		return err
	}

	header, items := sf.dataset.GetAllRows()

	// get indexes of items we are interested in for calculation
	quotaInx, err := findLocation(header, "QUOTA")
	if err != nil {
		return err
	}
	weekInx, err := findLocation(header, "WEEK")
	if err != nil {
		return err
	}
	w1yrInx, err := findLocation(header, "W1YR")
	if err != nil {
		return err
	}
	qrtrInx, err := findLocation(header, "QRTR")
	if err != nil {
		return err
	}
	addrInx, err := findLocation(header, "ADD")
	if err != nil {
		return err
	}
	wavfndInx, err := findLocation(header, "WAVFND")
	if err != nil {
		return err
	}
	hhldInx, err := findLocation(header, "HHLD")
	if err != nil {
		return err
	}
	persnoInx, err := findLocation(header, "PERSON")
	if err != nil {
		return err
	}

	for i := range column.Rows {
		row := items[i]

		quota, err := strconv.ParseFloat(row[quotaInx], 64)
		if err != nil {
			return err
		}

		week, err := strconv.ParseFloat(row[weekInx], 64)
		if err != nil {
			return err
		}

		w1yr, err := strconv.ParseFloat(row[w1yrInx], 64)
		if err != nil {
			return err
		}

		qrtr, err := strconv.ParseFloat(row[qrtrInx], 64)
		if err != nil {
			return err
		}

		addr, err := strconv.ParseFloat(row[addrInx], 64)
		if err != nil {
			return err
		}

		wavfnd, err := strconv.ParseFloat(row[wavfndInx], 64)
		if err != nil {
			return err
		}

		hhld, err := strconv.ParseFloat(row[hhldInx], 64)
		if err != nil {
			return err
		}

		persno, err := strconv.ParseFloat(row[persnoInx], 64)
		if err != nil {
			return err
		}

		n := (quota * 100000000000) + (week * 1000000000) + (w1yr * 100000000) +
			(qrtr * 10000000) + (addr * 100000) + (wavfnd * 10000) + (hhld * 100) + persno
		column.Rows[i] = int64(n)
	}

	return nil
}
