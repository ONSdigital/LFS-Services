package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"services/config"
	"services/db"
	"services/types"
	"testing"
)

// TODO: Ask Paul how to condense these types
type successfulTestCase struct {
	description      string
	year             string
	month            string
	expectedResponse types.OkayResponse
}

type xFailTestCase struct {
	description      string
	year             string
	month            string
	expectedResponse types.ErrorResponse
}

func cleanseMonthlyBatchTable(t *testing.T) {
	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gbBatchTable := config.Config.Database.GbBatchTable
	niBatchTable := config.Config.Database.NiBatchTable
	batchTable := config.Config.Database.MonthlyBatchTable

	tables := []string{gbBatchTable, niBatchTable, batchTable}

	// For each table: verify configuration set and then cleanse
	for _, table := range tables {
		if table == "" {
			t.Fatalf("%s table configuration not set", table)
		}
		if err := dbase.DeleteFrom(table); err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func TestMonth(t *testing.T) {
	// Cleansing and any setup
	cleanseMonthlyBatchTable(t)

	// TODO: Ask Paul how to prettify this!
	testCases := []successfulTestCase{
		{
			description:      "Create successful Monthly batch for May 2014",
			year:             "2014",
			month:            "05",
			expectedResponse: types.OkayResponse{Status: OK},
		},
		{
			description:      "Create successful Monthly batch for Jan 2016",
			year:             "2016",
			month:            "1",
			expectedResponse: types.OkayResponse{Status: OK},
		},
		{
			description:      "Create successful Monthly batch for Dec 2018",
			year:             "2018",
			month:            "12",
			expectedResponse: types.OkayResponse{Status: OK},
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest("POST", "/batches/monthly/", nil)
		r = mux.SetURLVars(r, map[string]string{"year": tc.year, "month": tc.month})
		w := httptest.NewRecorder()

		BatchHandler{}.CreateMonthlyBatchHandler(w, r)

		expectedResponse := tc.expectedResponse
		x, _ := json.Marshal(expectedResponse)
		expected := fmt.Sprintf("%s\n", string(x))

		if assert.Equal(t, http.StatusOK, w.Code) {
			if expected != w.Body.String() {
				t.Errorf("\n%s failed: \nExpected: %s\nActual: %s", tc.description, expected, w.Body.String())
			}

			t.Logf("\nPASS: %s", tc.description)
		}
	}
}

func TestMonthXFail(t *testing.T) {
	testCases := []xFailTestCase{
		{
			description: "Assert monthly batch already exists for Dec 2018",
			year:        "2018",
			month:       "12",
			expectedResponse: types.ErrorResponse{
				Status:       Error,
				ErrorMessage: "monthly batch for month 12, year 2018 already exists",
			},
		},
		{
			description: "Assert error for month 44",
			year:        "2014",
			month:       "44",
			expectedResponse: types.ErrorResponse{
				Status:       Error,
				ErrorMessage: "the month value is 44, must be between 1 and 12"},
		},
		{
			description: "Assert error for month Q",
			year:        "2018",
			month:       "Q",
			expectedResponse: types.ErrorResponse{
				Status:       Error,
				ErrorMessage: "invalid period: Q, expected one of 1-12",
			},
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest("POST", "/batches/monthly/", nil)
		r = mux.SetURLVars(r, map[string]string{"year": tc.year, "month": tc.month})
		w := httptest.NewRecorder()

		BatchHandler{}.CreateMonthlyBatchHandler(w, r)

		expectedResponse := tc.expectedResponse
		x, _ := json.Marshal(expectedResponse)
		expected := fmt.Sprintf("%s\n", string(x))

		if assert.Equal(t, http.StatusOK, w.Code) {
			if expected != w.Body.String() {
				t.Errorf("\n%s failed: \nExpected: %s\nActual: %s", tc.description, expected, w.Body.String())
			}
			// TODO: Improve output :)
			t.Log(w.Body.String())
			t.Logf("\nPASS: %s", tc.description)
		}
	}
}
