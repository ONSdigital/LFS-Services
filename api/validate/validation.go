package validate

import (
	"fmt"
	"math"
	"strconv"
)

type ValidationResult int

const ValidationFailed ValidationResult = 0
const ValidationSuccessful ValidationResult = 1

type Val struct {
	ValidationResponse
	error
}

type ValidationResponse struct {
	ValidationResult ValidationResult
	ErrorMessage     string
}

type Validation interface {
	Validate(period, year int) (ValidationResponse, error)
}

/*
Base validation. To use this, use composition in concrete structs
*/
type Validator struct {
	headers []string
	rows    [][]string
}

func (v Validator) GetRowsAsDouble(colName string) ([]float64, error) {
	var res []float64
	headers := v.headers
	rows := v.rows

	findRow := func() (int, bool) {
		for i, col := range headers {
			if col == colName {
				return i, true
			}
		}
		return 0, false
	}

	a, ok := findRow()
	if !ok {
		return nil, fmt.Errorf("cannot find column %s", colName)
	}

	for _, b := range rows {
		elem := b[a]
		val, err := strconv.ParseFloat(elem, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}

	return res, nil
}

func (v Validator) GetRowsAsString(colName string) ([]string, error) {
	var res []string

	findRow := func() (int, bool) {
		for i, col := range v.headers {
			if col == colName {
				return i, true
			}
		}
		return 0, false
	}

	a, ok := findRow()
	if !ok {
		return nil, fmt.Errorf("cannot find column %s", colName)
	}

	for _, b := range v.rows {
		res = append(res, b[a])
	}

	return res, nil
}

/*
Check if any rows in the list of columns to check are 'missing' where missing is defined as a NaN
for int and float types respectively.
*/
func (v Validator) validateMissingValues(columnsToCheck []string) (ValidationResponse, error) {
	for _, col := range columnsToCheck {

		floatCheck := func() (ValidationResponse, error) {
			rows, err := v.GetRowsAsDouble(col)
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
			rows, err := v.GetRowsAsString(col)
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

		if col == "PCODE" {
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
