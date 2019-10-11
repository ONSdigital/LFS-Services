package validate

import "services/dataset"

type ValidationResult int

const ValidationFailed ValidationResult = 0
const ValidationSuccessful ValidationResult = 1

type ValidationResponse struct {
	ValidationResult ValidationResult
	ErrorMessage     string
}

/**
Base validation. To use this, use composition in concrete structs
*/
type Validator struct {
	dataset *dataset.Dataset
}
