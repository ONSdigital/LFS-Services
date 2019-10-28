package api

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
