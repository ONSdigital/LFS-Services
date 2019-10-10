package validation

const (
	ValidationFailed = iota
	ValidationSuccessful
)

type ValidationResponse struct {
	ValidationStatus int
	ErrorMessage     string
}

type Validation interface {
	Validate() (ValidationResponse, error)
}
