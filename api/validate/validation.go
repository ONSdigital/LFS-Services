package validate

type ValidationResult int

const ValidationFailed ValidationResult = 0
const ValidationSuccessful ValidationResult = 1

type ValidationResponse struct {
	ValidationResult ValidationResult
	ErrorMessage     string
}

type Validator interface {
	Validate() (ValidationResponse, error)
}
