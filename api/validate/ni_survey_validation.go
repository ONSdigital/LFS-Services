package validate

import (
	"fmt"
	"services/types"
	"time"
)

type NISurveyValidation struct {
	Validator
}

func NewNISurveyValidation(data *types.SavImportData) NISurveyValidation {
	return NISurveyValidation{Validator: Validator{data}}
}

func (sf NISurveyValidation) Validate(period, year int) (ValidationResponse, error) {

	var columnsToCheck = []string{
		"REFDTE", "PCODE", "QUOTA", "WEEK", "W1YR", "QRTR", "ADDR", "WAVFND", "HHLD", "PERSNO",
	}

	v, e := sf.validateMissingValues(columnsToCheck)

	if e != nil {
		return v, e
	}

	v, e = sf.ValidateNIDates(period, year)
	return v, e

}

func (sf NISurveyValidation) ValidateNIDates(month, year int) (ValidationResponse, error) {

	// get the list of weeks in the sav file
	rows, err := sf.GetRowsAsDouble("REFDTE")
	if err != nil {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     err.Error(),
		}, err
	}

	var weeks = make(map[float64]int, 0)
	for pos, j := range rows {
		weeks[j] = pos
	}

	// check how many weeks we have
	if len(weeks) != 4 && len(weeks) != 5 {
		return ValidationResponse{
			ValidationResult: ValidationFailed,
			ErrorMessage:     fmt.Sprintf("rows must contain either 4 or 5 weeks of data, found: %d", len(weeks)),
		}, fmt.Errorf("rows must contain either 4 or 5 weeks of data found: %d", len(weeks))
	}

	//check each week starts on a sunday
	for timeStamp, _ := range weeks {
		i := int64(timeStamp) - (141428 * 86400)
		tm := time.Unix(i, 0)
		if tm.Weekday() != 0 {
			return ValidationResponse{
				ValidationResult: ValidationFailed,
				ErrorMessage:     fmt.Sprintf("RFEDTE is not a Sunday - it is a %s", tm.Weekday().String()),
			}, fmt.Errorf(fmt.Sprintf("RFEDTE is not a Sunday - it is a %s", tm.Weekday().String()))
		}

		// check week number against month and year
		m := int(tm.Month())
		y := tm.Year()

		if m != month {
			return ValidationResponse{
				ValidationResult: ValidationFailed,
				ErrorMessage:     fmt.Sprintf("Week number in RFEDTE is not the required month %d, it is %d", month, m),
			}, fmt.Errorf(fmt.Sprintf("week number in RFEDTE is not the required month %d, it is %d", month, y))
		}

		if y != year {
			return ValidationResponse{
				ValidationResult: ValidationFailed,
				ErrorMessage:     fmt.Sprintf("Year number in RFEDTE is not the required year %d, it is %d", year, y),
			}, fmt.Errorf(fmt.Sprintf("year number in RFEDTE is not the required year %d, it is %d", year, y))
		}

	}

	return ValidationResponse{
		ValidationResult: ValidationSuccessful,
		ErrorMessage:     "",
	}, nil
}
