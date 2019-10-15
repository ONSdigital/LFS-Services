package api

import (
	"fmt"
	"gopkg.in/validator.v2"
	_ "services/api/validate"
)

type UserCredentials struct {
	Username string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

func (h RestHandlers) login(username string, password string) error {
	// Validate user input
	userCreds := UserCredentials{Username: username, Password: password}

	if errs := validator.Validate(userCreds); errs != nil {
		return fmt.Errorf("Invalid Username or Password.")
	}

	//TODO: Get user creds from database

	// Winner, winner, chicken dinner!
	return nil
}
