package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"services/config"
	"services/db"
	"services/types"
	"testing"
)

func cleanseMonthlyBatchTable() error {
	gbBatchTable := config.Config.Database.GbBatchTable
	if gbBatchTable == "" {
		panic("gb batch table configuration not set")
	}

	niBatchTable := config.Config.Database.NiBatchTable
	if niBatchTable == "" {
		panic("gb batch table configuration not set")
	}

	batchTable := config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		panic("monthly batch table configuration not set")
	}

	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		return err
	}

	// First cleanse gb_batch_items due to constraints on foreign key
	if err = dbase.DeleteFrom(gbBatchTable); err != nil {
		return err
	}

	// Then cleanse ni_batch_items due to constraints on foreign key
	if err = dbase.DeleteFrom(niBatchTable); err != nil {
		return err
	}

	// Finally cleanse monthly_batch table
	if err = dbase.DeleteFrom(batchTable); err != nil {
		return err
	}

	return nil
}

func createMonthlyBatchSetup(cleanse bool, vars map[string]string) (bool, *httptest.ResponseRecorder) {
	if cleanse {
		// Cleanse monthly tables
		err := cleanseMonthlyBatchTable()

		if err != nil {
			log.Error().Msg(err.Error())
			return false, nil
		}
	}

	year := vars["year"]
	month := vars["month"]

	url := fmt.Sprintf("/batches/monthly/%s/%s", year, month)

	r, _ := http.NewRequest("POST", url, nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, vars)

	BatchHandler{}.CreateMonthlyBatchHandler(w, r)

	return true, w

}

// Test to assert successful creation of monthly batch for May 2014
func TestCreateMonthlyBatchMay2014(t *testing.T) {
	vars := map[string]string{
		"year":  "2014",
		"month": "06",
	}
	expectedResponse := types.OkayResponse{Status: OK}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	res, w := createMonthlyBatchSetup(true, vars)

	if res == false {
		log.Error().Msg("createMonthlyBatchSetup() didn't return anything :(")
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", string(expected)), w.Body.String())
}

// Test to assert successful creation of monthly batch for Jan 2016
func TestCreateMonthlyBatchJan2016(t *testing.T) {
	vars := map[string]string{
		"year":  "2016",
		"month": "1",
	}
	expectedResponse := types.OkayResponse{Status: OK}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	res, w := createMonthlyBatchSetup(true, vars)

	if res == false {
		log.Error().Msg("createMonthlyBatchSetup() didn't return anything :(")
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", string(expected)), w.Body.String())
}

// Test to assert successful creation of monthly batch for Dec 2018
func TestCreateMonthlyBatchDec2018(t *testing.T) {
	vars := map[string]string{
		"year":  "2018",
		"month": "12",
	}
	expectedResponse := types.OkayResponse{Status: OK}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	res, w := createMonthlyBatchSetup(true, vars)

	if res == false {
		log.Error().Msg("createMonthlyBatchSetup() didn't return anything :(")
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", string(expected)), w.Body.String())
}

// Test to assert monthly batch already exists for Dec 2018
func TestCreateMonthlyBatchMay2014Fail(t *testing.T) {
	vars := map[string]string{
		"year":  "2018",
		"month": "12",
	}
	//expectedResponse := types.OkayResponse{Status:OK}
	expectedResponse := types.ErrorResponse{
		Status:       Error,
		ErrorMessage: "monthly batch for month 12, year 2018 already exists",
	}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	res, w := createMonthlyBatchSetup(false, vars)

	if res == false {
		log.Error().Msg("createMonthlyBatchSetup() didn't return anything :(")
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", string(expected)), w.Body.String())
}

// Test to assert error for month 44
func TestCreateMonthlyBatchBadMonthFail(t *testing.T) {
	vars := map[string]string{
		"year":  "2018",
		"month": "44",
	}
	//expectedResponse := types.OkayResponse{Status:OK}
	expectedResponse := types.ErrorResponse{
		Status:       Error,
		ErrorMessage: "the month value is 44, must be between 1 and 12",
	}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	res, w := createMonthlyBatchSetup(false, vars)

	if res == false {
		log.Error().Msg("createMonthlyBatchSetup() didn't return anything :(")
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", string(expected)), w.Body.String())
}

// Test to assert error for month Q
func TestCreateMonthlyBatchStringMonthFail(t *testing.T) {
	vars := map[string]string{
		"year":  "2018",
		"month": "Q",
	}
	//expectedResponse := types.OkayResponse{Status:OK}
	expectedResponse := types.ErrorResponse{
		Status:       Error,
		ErrorMessage: "invalid period: Q, expected one of 1-12",
	}
	expected, err := json.Marshal(expectedResponse)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	res, w := createMonthlyBatchSetup(false, vars)

	if res == false {
		log.Error().Msg("createMonthlyBatchSetup() didn't return anything :(")
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("%s\n", string(expected)), w.Body.String())
}
