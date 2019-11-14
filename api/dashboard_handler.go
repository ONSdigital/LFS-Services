package api

import (
	"fmt"
	"net/http"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (d DashboardHandler) HandleDashboardRequest(w http.ResponseWriter, r *http.Request) {

	//Functionality
	res, err := d.GetDashboardInfo()

	// Error handling
	if err != nil {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: err.Error()}.sendResponse(w, r)
		return
	}

	if len(res) == 0 {
		ErrorResponse{
			Status:       Error,
			ErrorMessage: fmt.Sprintf("There is nothing to show")}.sendResponse(w, r)
		return
	}

	SendDataResponse{}.sendResponse(w, r, res)

}
