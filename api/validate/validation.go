package validate

type ValidationResult int

const ValidationFailed ValidationResult = 0
const ValidationSuccessful ValidationResult = 1

type ValidationResponse struct {
	ValidationResult ValidationResult
	ErrorMessage     string
}

/*
Base validation. To use this, use composition in concrete structs
*/
type Validator struct {
	Headers *[]string
	Rows    *[][]string
}
