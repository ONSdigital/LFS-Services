package tests

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"services/api"
	"services/config"
	"services/db"
	"testing"
)

func CreateMonthlyBatchSetup(cleanse bool, vars map[string]string) (bool, *httptest.ResponseRecorder) {
	if cleanse {
		// Cleanse monthly tables
		err := CleanseMonthlyBatchTable()

		if err != nil {
			log.Error().Msg("Couldn't clear tables. You've made El sad :'( I refuse to continue testing")
			return false, nil
		}
	}

	year := vars["year"]
	month := vars["month"]

	url := fmt.Sprintf("/batches/monthly/%s/%s", year, month)

	r, _ := http.NewRequest("POST", url, nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, vars)

	api.BatchHandler{}.CreateMonthlyBatchHandler(w, r)

	return true, w

}

func CleanseMonthlyBatchTable() error {
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

	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
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

// Test to assert successful creation of monthly batch for May 2014
func TestCreateMonthlyBatchMay2014(t *testing.T) {
	vars := map[string]string{
		"year":  "2014",
		"month": "05",
	}

	res, w := CreateMonthlyBatchSetup(true, vars)

	if res == false {
		a := assert.AnError
		print(a)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// TODO: In need of a Paul-Soule-special-tutorial on strings
	expected := `{"status":"OK"}`
	assert.Equal(t, expected, w.Body.String())
}

// Test to assert successful creation of monthly batch for Jan 2016
func TestCreateMonthlyBatchJan2016(t *testing.T) {
	vars := map[string]string{
		"year":  "2016",
		"month": "1",
	}

	res, w := CreateMonthlyBatchSetup(true, vars)

	if res == false {
		a := assert.AnError
		print(a)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// TODO: In need of a Paul-Soule-special-tutorial on strings
	expected := `{"status":"OK"}`
	assert.Equal(t, expected, w.Body.String())
}

// Test to assert successful creation of monthly batch for Dec 2018
func TestCreateMonthlyBatchDec2018(t *testing.T) {
	vars := map[string]string{
		"year":  "2018",
		"month": "12",
	}

	res, w := CreateMonthlyBatchSetup(true, vars)

	if res == false {
		a := assert.AnError
		print(a)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// TODO: In need of a Paul-Soule-special-tutorial on strings
	expected := `{"status":"OK"}`
	assert.Equal(t, expected, w.Body.String())
}

// Test to assert monthly batch already exists for May 2014
func TestCreateMonthlyBatchMay2014Fail(t *testing.T) {
	vars := map[string]string{
		"year":  "2014",
		"month": "05",
	}

	res, w := CreateMonthlyBatchSetup(false, vars)

	if res == false {
		a := assert.AnError
		print(a)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// TODO: In need of a Paul-Soule-special-tutorial on strings
	expected := `{"status":"ERROR","errorMessage":"monthly batch for month 5, year 2014 already exists"}
` // Don't delete this line, I swear down!
	assert.Equal(t, expected, w.Body.String())
}
