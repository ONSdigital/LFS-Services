package api

/*
input file names / types
*/
const (
	SurveyFile  = "Survey"
	AddressFile = "Address"
)

const (
	Error = "ERROR"
	OK    = "OK"
)

type RestHandlers struct{}

/*
Create a new RestHandler
*/
func NewRestHandler() *RestHandlers {
	return &RestHandlers{}
}
