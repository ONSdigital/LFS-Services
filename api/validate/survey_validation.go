package validate

import (
	"fmt"
	"math"
	"services/dataset"
	"time"
)

type SurveyValidation struct {
	Validator
}

func NewSurveyValidation(dataset *dataset.Dataset) SurveyValidation {
	return SurveyValidation{Validator: Validator{dataset}}
}

type Val struct {
	ValidationResponse
	error
}

func (sf SurveyValidation) Validate() (ValidationResponse, error) {

	v, e := sf.validateMissingValues()

	if e != nil {
		return v, e
	}

	v, e = sf.validateREFDTE()
	return v, e

}

var columnsToCheck = []string{"REFDTE", "PCODE", "QUOTA", "WEEK", "W1YR", "QRTR", "ADD", "WAVFND", "HHLD", "PERSNO"}

/*
Check if any rows in the list of columns to check are 'missing' where missing is defined as a NaN
for int and float types respectively.
*/
func (sf SurveyValidation) validateMissingValues() (ValidationResponse, error) {
	for _, v := range columnsToCheck {

		floatCheck := func() (ValidationResponse, error) {
			rows, err := sf.dataset.GetRowsAsDouble(v)
			if err != nil {
				return ValidationResponse{
					ValidationResult: ValidationFailed,
					ErrorMessage:     err.Error(),
				}, err
			}
			for _, j := range rows {
				if math.IsNaN(j) {
					return ValidationResponse{
						ValidationResult: ValidationFailed,
						ErrorMessage:     "column %s has a missing value",
					}, fmt.Errorf("column %s has a missing value - NaN", v)
				}
			}
			return ValidationResponse{
				ValidationResult: ValidationSuccessful,
				ErrorMessage:     "Successful",
			}, nil
		}

		stringCheck := func() (ValidationResponse, error) {
			rows, err := sf.dataset.GetRowsAsString(v)
			if err != nil {
				return ValidationResponse{
					ValidationResult: ValidationFailed,
					ErrorMessage:     err.Error(),
				}, err
			}
			for _, j := range rows {
				if j == "" {
					return ValidationResponse{
						ValidationResult: ValidationFailed,
						ErrorMessage:     "column %s has a missing value",
					}, fmt.Errorf("column %s has a missing value", v)
				}
			}
			return ValidationResponse{
				ValidationResult: ValidationSuccessful,
				ErrorMessage:     "Successful",
			}, nil
		}

		if v == "PCODE" {
			return stringCheck()
		} else {
			return floatCheck()
		}

	}

	return ValidationResponse{
		ValidationResult: ValidationSuccessful,
		ErrorMessage:     "Successful",
	}, nil
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
func (sf SurveyValidation) validateREFDTE() (ValidationResponse, error) {
	rows, err := sf.dataset.GetRowsAsDouble("REFDTE")
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
				ErrorMessage:     "rows contain different values for RFEDTE",
			}, fmt.Errorf("rows contain different values for RFEDTE")
		}
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

	return ValidationResponse{
		ValidationResult: ValidationSuccessful,
		ErrorMessage:     "",
	}, nil
}
