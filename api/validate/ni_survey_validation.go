package validate

import (
	"services/types"
)

type NISurveyValidation struct {
	Validator
}

func NewNISurveyValidation(data [][]string) NISurveyValidation {
	return NISurveyValidation{Validator: Validator{data}}
}

func (sf NISurveyValidation) Validate(period, year int) (ValidationResponse, error) {

	var columnsToCheck = []string{
		"RefDte", "PCode", "Quota", "Week", "W1Yr", "Qrtr", "Addr", "WavFnd", "Hhld", "PersNo",
	}

	v, e := sf.validateMissingValues(columnsToCheck)

	if e != nil {
		return v, e
	}

	v, e = sf.validateREFDTE(period, year, types.NI)
	return v, e

}
