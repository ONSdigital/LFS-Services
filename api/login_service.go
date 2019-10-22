package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
	_ "services/api/validate"
	"services/db"
	"services/types"
	"strings"
)

func (h RestHandlers) login(username string, password string) error {
	log.Debug().
		Str("client", h.r.RemoteAddr).
		Str("uri", h.r.RequestURI).
		Msg("Validating login input.")

	// Validate user input
	userCreds := types.UserCredentials{Username: username, Password: password}

	if errs := validator.Validate(userCreds); errs != nil {
		log.Error().
			Msg("Invalid Username or Password.")
		return fmt.Errorf("Invalid Username or Password.")
	}

	log.Debug().
		Str("client", h.r.RemoteAddr).
		Str("uri", h.r.RequestURI).
		Msg("Retrieving user credentials from database.")

	// Get user creds from database
	creds, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Msg(err.Error())
		return err
	}
	user, err := creds.GetUserID(username)
	if err != nil {
		return err
	}

	log.Debug().
		Str("client", h.r.RemoteAddr).
		Str("uri", h.r.RequestURI).
		Msg("Assert user credentials match.")

	// Compare and assert credentials match
	matchErr := comparePasswords(user.Password, password)

	if strings.Compare(username, user.Username) != 0 || matchErr == false {
		log.Error().
			Msg("Invalid username or password")
		return fmt.Errorf("Invalid username or password")
	}

	return nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false
	}
	return true
}
