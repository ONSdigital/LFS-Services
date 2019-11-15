package validate

import (
	"fmt"
	"math"
	"services/types"
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

func (v Validator) findRowIndex(colName string) (int, bool) {
	for i, col := range v.data.Header {
		if col.VariableName == colName {
			return i, true
		}
	}
	return 0, false
}

/*
Base validation. To use this, use composition in concrete structs
*/
type Validator struct {
	data *types.SavImportData
}

func (v Validator) GetRowsAsDouble(colName string) ([]float64, error) {
	var res []float64

	a, ok := v.findRowIndex(colName)
	if !ok {
		return nil, fmt.Errorf("cannot find column %s", colName)
	}

	for _, b := range v.data.Rows {
		elem := b.RowData[a]
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

	a, ok := v.findRowIndex(colName)
	if !ok {
		return nil, fmt.Errorf("cannot find column %s", colName)
	}

	for _, b := range v.data.Rows {
		res = append(res, b.RowData[a])
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
