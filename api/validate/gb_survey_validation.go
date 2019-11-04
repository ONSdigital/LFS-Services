package validate

import (
	"services/types"
)

type GBSurveyType int

type GBSurveyValidation struct {
	Validator
}

func NewGBSurveyValidation(data [][]string) GBSurveyValidation {
	return GBSurveyValidation{Validator: Validator{data}}
}

func (sf GBSurveyValidation) Validate(period, year int) (ValidationResponse, error) {

	var columnsToCheck = []string{
		"RefDte", "PCode", "Quota", "Week", "W1Yr", "Qrtr", "Addr", "WavFnd", "Hhld", "PersNo",
	}
	v, e := sf.validateMissingValues(columnsToCheck)

	if e != nil {
		return v, e
	}

	v, e = sf.validateREFDTE(period, year, types.GB)
	return v, e

}
