package validate

import (
	"fmt"
	"services/types"
	"time"
)

type GBSurveyType int

type GBSurveyValidation struct {
	Validator
}

func NewGBSurveyValidation(headers []string, rows [][]string, data types.SavImportData) GBSurveyValidation {
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

	v, e = sf.validateREFDTE(period, year)
	return v, e

}

/*
Validate the REFDTE field in the Survey file

SPSS stores timestamps as the numbers of seconds between the year 1582 (start of the Gregorian calendar)
and a given time on a given date. To get the actual date from this we need to:

1. Get the difference between the Gregorian time and the Unix epoch in seconds (141428)
2. Multiply this value by the number of seconds in a day (86400)
3. Subtract this value from the SPSS timestamp to get the Unix time, and
4. Get the date from the Unix time using standard Go functions.

*/
func (v Validator) validateREFDTE(period, year int) (ValidationResponse, error) {
	rows, err := v.GetRowsAsDouble("REFDTE")
	if err != nil {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     err.Error(),
		}, err
	}

	var val float64
	for _, j := range rows {
		if val == 0.0 {
			val = j
		}
		if val != j {
			return ValidationResponse{
				ValidationResult: ValidationFailed,
				ErrorMessage:     "rows contain different values for REFDTE",
			}, fmt.Errorf("rows contain different values for REFDTE")
		}
	}

	if len(rows) == 0 {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     fmt.Sprintf("There are no rows to validate"),
		}, fmt.Errorf(fmt.Sprintf("there are no rows to validate"))
	}

	// Take the first row rather than checking in a loop
	i := int64(rows[0]) - (141428 * 86400)
	tm := time.Unix(i, 0)
	if tm.Weekday() != 0 {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     fmt.Sprintf("RFEDTE is not a Sunday - it is a %s", tm.Weekday().String()),
		}, fmt.Errorf(fmt.Sprintf("RFEDTE is not a Sunday - it is a %s", tm.Weekday().String()))
	}

	// check week number against RFEDTE
	y, w := tm.ISOWeek()

	if w != period {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     fmt.Sprintf("Week number in RFEDTE is not the required week %d, it is %d", period, w),
		}, fmt.Errorf(fmt.Sprintf("week number in RFEDTE is not the required week %d, it is %d", period, w))
	}

	if y != year {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     fmt.Sprintf("Year number in RFEDTE is not the required year %d, it is %d", year, y),
		}, fmt.Errorf(fmt.Sprintf("year number in RFEDTE is not the required year %d, it is %d", year, y))
	}

	return ValidationResponse{
		ValidationResult: ValidationSuccessful,
		ErrorMessage:     "",
	}, nil
}
