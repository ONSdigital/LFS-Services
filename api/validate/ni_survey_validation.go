package validate

import (
	"services/types"
)

type NISurveyValidation struct {
	Validator
}

func NewNISurveyValidation(headers []string, rows [][]string) NISurveyValidation {
	return NISurveyValidation{Validator: Validator{headers, rows}}
}

func (sf NISurveyValidation) Validate(period, year int) (ValidationResponse, error) {

	var columnsToCheck = []string{
		"REFDTE", "PCODE", "QUOTA", "WEEK", "W1YR", "QRTR", "ADDR", "WAVFND", "HHLD", "PERSNO",
	}

	v, e := sf.validateMissingValues(columnsToCheck)

	if e != nil {
		return v, e
	}

	v, e = sf.validateREFDTE(period, year, types.NI)
	return v, e

}
