package validate

import (
	"services/types"
)

type GBSurveyType int

type GBSurveyValidation struct {
	Validator
}

func NewGBSurveyValidation(headers []string, rows [][]string) GBSurveyValidation {
	return GBSurveyValidation{Validator: Validator{headers, rows}}
}

func (sf GBSurveyValidation) Validate(period, year int) (ValidationResponse, error) {

	var columnsToCheck = []string{
		"REFDTE", "PCODE", "QUOTA", "WEEK", "W1YR", "QRTR", "ADDR", "WAVFND", "HHLD", "PERSNO",
	}
	v, e := sf.validateMissingValues(columnsToCheck)

	if e != nil {
		return v, e
	}

	v, e = sf.validateREFDTE(period, year, types.GB)
	return v, e

}
