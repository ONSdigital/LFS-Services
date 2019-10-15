package api

import "fmt"

func (h RestHandlers) login(username string) error {
	//TODO: Get user creds from database

	return fmt.Errorf("%s is a tit", username)
}
