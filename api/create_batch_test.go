package api

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"services/config"
	"services/db"
	"services/types"
	"testing"
)

type testCase struct {
	description  string
	year         string
	period       string
	expectedCode int
}

func TestMonthlyBatchHandler(t *testing.T) {
	// Cleansing and any setup
	tearDown(t)

	r := httptest.NewRequest("POST", "/batches/monthly/", nil)

	for _, mtc := range monthlyTestCases() {
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, map[string]string{"year": mtc.year, "month": mtc.period})

		BatchHandler{}.CreateMonthlyBatchHandler(w, r)

		if !assert.Equal(t, mtc.expectedCode, w.Code) {
			t.Fatalf("\nFAILED TEST CASE: %s\nERROR MESSAGE: %s",
				mtc.description, w.Body.String())
		}
		t.Logf("\nPASSED: %s\nPERIOD: %s, %s\nBODY: %s",
			mtc.description, mtc.year, mtc.period, w.Body.String())
	}
}

func TestQuarterlyBatchHandler(t *testing.T) {
	// Cleansing and any setup
	tearDown(t)
	setupMonthlyTables(t, 1, 2, 2015, 0)  // Jan-Feb 2015 status 0
	setupMonthlyTables(t, 1, 2, 2016, 4)  // Jan-Feb 2016 status 4
	setupMonthlyTables(t, 3, 1, 2016, 0)  // Mar 2016 status 0
	setupMonthlyTables(t, 1, 12, 2017, 4) // Jan-Dec 2017 status 4

	r := httptest.NewRequest("POST", "/batches/quarterly/", nil)

	for _, qtc := range quarterlyTestCases() {
		r = mux.SetURLVars(r, map[string]string{"year": qtc.year, "quarter": qtc.period})
		w := httptest.NewRecorder()

		BatchHandler{}.CreateQuarterlyBatchHandler(w, r)

		if !assert.Equal(t, qtc.expectedCode, w.Code) {
			t.Fatalf("\n\nFAILED TEST CASE: %s\nPERIOD: %s, %s\nERROR MESSAGE: %s",
				qtc.description, qtc.period, qtc.year, w.Body.String())
		}
		t.Logf("\nPASSED: %s\nPERIOD: %s, %s\nBODY: %s",
			qtc.description, qtc.year, qtc.period, w.Body.String())
	}
}

func TestAnnualBatchHandler(t *testing.T) {
	// Cleansing and any setup
	tearDown(t)

	setupMonthlyTables(t, 1, 6, 2015, 0)  // Jan-Jun 2015 status 0
	setupMonthlyTables(t, 1, 9, 2016, 4)  // Jan-Sep 2016 status 4
	setupMonthlyTables(t, 10, 2, 2016, 0) // Oct-Dec 2016 status 0
	setupMonthlyTables(t, 1, 12, 2017, 4) // Jan-Dec 2017 status 4

	setupQuarterlyTables(t, 1, 2, 2015, 0) // Q1-Q2 2015 status 0
	setupQuarterlyTables(t, 1, 3, 2016, 4) // Q1-Q3 2016 status 4
	setupQuarterlyTables(t, 4, 1, 2016, 0) // Q4 2016 status 0
	setupQuarterlyTables(t, 1, 4, 2017, 4) // Q1-Q4 2017 status 4

	r := httptest.NewRequest("POST", "/batches/annual/", nil)

	for _, atc := range annualTestCases() {
		r = mux.SetURLVars(r, map[string]string{"year": atc.year})
		w := httptest.NewRecorder()

		BatchHandler{}.CreateAnnualBatchHandler(w, r)

		if !assert.Equal(t, atc.expectedCode, w.Code) {
			t.Fatalf("\n\nFAILED TEST CASE: %s\nERROR MESSAGE: %s",
				atc.description, w.Body.String())
		}
		t.Logf("\nPASSED: %s\nYEAR %s\nBODY: %s",
			atc.description, atc.year, w.Body.String())
	}
}

func tearDown(t *testing.T) {
	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gbBatchTable := config.Config.Database.GbBatchTable
	niBatchTable := config.Config.Database.NiBatchTable
	batchTable := config.Config.Database.MonthlyBatchTable
	quarterlyBatchTable := config.Config.Database.QuarterlyBatchTable
	annualBatchTable := config.Config.Database.AnnualBatchTable

	tables := []string{gbBatchTable, niBatchTable, batchTable, quarterlyBatchTable, annualBatchTable}

	// For each table: confirm configuration is set and then cleanse
	for _, table := range tables {
		if table == "" {
			t.Fatal("table configuration not set")
		}
		if err := dbase.DeleteFrom(table); err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func setupMonthlyTables(t *testing.T, month, count, year, status int) {
	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		t.Fatalf(err.Error())
	}

	batchTable := config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		t.Fatal("monthly_batch table configuration not set")
	}

	// Insert a load of mock data and set status to 4
	for i := month; i <= count; i++ {
		batch := types.MonthlyBatch{
			Year:        year,
			Month:       i,
			Status:      status,
			Description: "Mock data for Testing",
		}
		if err := dbase.CreateMonthlyBatch(batch); err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func setupQuarterlyTables(t *testing.T, quarter, count, year, status int) {
	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		t.Fatalf(err.Error())
	}

	quarterlyTable := config.Config.Database.QuarterlyBatchTable
	if quarterlyTable == "" {
		t.Fatal("quarterly_batch table configuration not set")
	}

	// Insert a load of mock data and set status to 4
	for i := quarter; i <= count; i++ {
		batch := types.QuarterlyBatch{
			Quarter:     i,
			Year:        year,
			Status:      status,
			Description: "Mock data for Testing",
		}
		if err := dbase.CreateQuarterlyBatch(batch); err != nil {
			t.Fatalf(err.Error())
		}
	}

}

func monthlyTestCases() []testCase {
	testCases := []testCase{
		{
			description:  "Monthly",
			year:         "2014",
			period:       "05",
			expectedCode: 200,
		},
		{
			description:  "Monthly",
			year:         "2016",
			period:       "1",
			expectedCode: 200,
		},
		{
			description:  "Monthly",
			year:         "2018",
			period:       "12",
			expectedCode: 200,
		},
		{
			description:  "Monthly (Assert Error: Already exists)",
			year:         "2018",
			period:       "12",
			expectedCode: 400,
		},
		{
			description:  "Monthly (Assert Error: Invalid month, expected month one of 1-12)",
			year:         "2014",
			period:       "44",
			expectedCode: 400,
		},
		{
			description:  "Monthly (Assert Error: Invalid month, expected month one of 1-12)",
			year:         "2018",
			period:       "Q",
			expectedCode: 400,
		},
		{
			description:  "Monthly (Assert Error: Expected year as integer)",
			year:         "Q",
			period:       "4",
			expectedCode: 400,
		},
	}

	return testCases
}

func quarterlyTestCases() []testCase {
	testCases := []testCase{
		{
			description:  "Quarterly",
			year:         "2017",
			period:       "Q4",
			expectedCode: 200,
		},
		{
			description:  "Quarterly",
			year:         "2017",
			period:       "Q3",
			expectedCode: 200,
		},
		{
			description:  "Quarterly",
			year:         "2017",
			period:       "Q2",
			expectedCode: 200,
		},
		{
			description:  "Quarterly",
			year:         "2017",
			period:       "Q1",
			expectedCode: 200,
		},
		{
			description:  "Quarterly (Assert Error: Already exists)",
			year:         "2017",
			period:       "Q4",
			expectedCode: 400,
		},
		{
			description:  "Quarterly (Assert Error: no monthly batches exist for Q4, 2019. Required 3 monthly batches to continue)",
			year:         "2014",
			period:       "Q4",
			expectedCode: 400,
		},
		{
			description:  "Quarterly (Assert Error: 2 monthly batches exist for Q1, 2015. Required 3 monthly batches to continue)",
			year:         "2015",
			period:       "Q1",
			expectedCode: 400,
		},
		{
			description:  "Quarterly (Assert Error: 2 valid monthly batches exist for Q1, 2016. Required 3 valid monthly batches to continue)",
			year:         "2016",
			period:       "Q1",
			expectedCode: 400,
		},
		{
			description:  "Quarterly (Assert Error: Invalid period, expected one of Q1-Q4)",
			year:         "2019",
			period:       "4",
			expectedCode: 400,
		},
		{
			description:  "Quarterly (Assert Error: Invalid period, expected one of Q1-Q4)",
			year:         "2019",
			period:       "Q5",
			expectedCode: 400,
		},
		{
			description:  "Quarterly (Assert Error: Invalid year, expected integer)",
			year:         "El is amazing",
			period:       "Q4",
			expectedCode: 400,
		},
	}

	return testCases
}

func annualTestCases() []testCase {
	testCases := []testCase{
		{
			description:  "Annual",
			year:         "2017",
			expectedCode: 200,
		},
		{
			description:  "Annual (Assert Error: Already exists)",
			year:         "2017",
			expectedCode: 400,
		},
		{
			description:  "Annual (Assert Error: 0 monthly batches for 2014. 12 months required to continue)",
			year:         "2014",
			expectedCode: 400,
		},
		{
			description:  "Annual (Assert Error: 12 months required)",
			year:         "2015",
			expectedCode: 400,
		},
		{
			description:  "Annual (Assert Error: 12 valid months required)",
			year:         "2016",
			expectedCode: 400,
		},
		{
			description:  "Annual (Assert Error: 4 quarter batches required)",
			year:         "2015",
			expectedCode: 400,
		},
		{
			description:  "Annual (Assert Error: 4 valid quarters required)",
			year:         "2016",
			expectedCode: 400,
		},
	}

	return testCases
}
