package validate

import (
	"fmt"
	"math"
	"services/dataset"
	"time"
)

const (
	GB = iota
	NI
)

type SurveyType int

type SurveyValidation struct {
	Validator
	Origin SurveyType
}

func NewSurveyValidation(dataset *dataset.Dataset, origin SurveyType) SurveyValidation {
	return SurveyValidation{Validator: Validator{dataset}, Origin: origin}
}

type Val struct {
	ValidationResponse
	error
}

func (sf SurveyValidation) Validate(period, year int) (ValidationResponse, error) {

	v, e := sf.validateMissingValues()

	if e != nil {
		return v, e
	}

	v, e = sf.validateREFDTE(period, year)
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
func (sf SurveyValidation) validateREFDTE(period, year int) (ValidationResponse, error) {
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

	if len(rows) == 0 {
		// No rows
		var typ = "GB"
		if sf.Origin == NI {
			typ = "NI"
		}
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     fmt.Sprintf("There are no rows to validate. Are you sure this is a %s file?", typ),
		}, fmt.Errorf(fmt.Sprintf("there are no rows to validate. Are you sure this is a %s file?", typ))
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

	if sf.Origin == GB {
		if w != period {
			return ValidationResponse{
				ValidationResult: ValidationFailed,
				ErrorMessage:     fmt.Sprintf("Week number in RFEDTE is not the required week %d, it is %d", period, w),
			}, fmt.Errorf(fmt.Sprintf("week number in RFEDTE is not the required week %d, it is %d", period, w))
		}
	} else {
		//NI
		if int(tm.Month()) != period {
			return ValidationResponse{
				ValidationResult: ValidationFailed,
				ErrorMessage:     fmt.Sprintf("Month number in RFEDTE is not the required month %d, it is %d", period, w),
			}, fmt.Errorf(fmt.Sprintf("month number in RFEDTE is not the required month %d, it is %d", period, w))
		}
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
