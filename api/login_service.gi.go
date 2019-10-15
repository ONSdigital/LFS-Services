package api

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
	_ "services/api/validate"
	"services/db"
	"services/types"
)

func (h RestHandlers) login(username string, password string) error {
	// Validate user input
	userCreds := types.UserCredentials{Username: username, Password: password}

	if errs := validator.Validate(userCreds); errs != nil {
		return fmt.Errorf("Invalid Username or Password.")
	}

	// Get user creds from database
	creds, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		return err
	}
	user, err := creds.GetUserID(username)
	if err != nil {
		return err
	}

	// Compare passwords
	matchErr := comparePasswords(user.Password, password)
	if matchErr != nil {
		return matchErr
	}

	// Winner, winner, chicken dinner!
	return nil
}

func comparePasswords(hashedPwd string, plainPwd string) error {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return err
	}
	return nil
}
