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
	year         string
	period       string
	expectedCode int
}

func TestMonthlyMay2014Success(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "05",
		expectedCode: 200,
	}
	tearDown(t)
	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestMonthlyJan2016Success(t *testing.T) {
	var tc = testCase{
		year:         "2016",
		period:       "1",
		expectedCode: 200,
	}
	tearDown(t)
	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestMonthlyDec2018Success(t *testing.T) {
	var tc = testCase{
		year:         "2018",
		period:       "12",
		expectedCode: 200,
	}
	tearDown(t)
	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestMonthlyAlreadyExistsXFail(t *testing.T) {
	var tc = testCase{
		year:         "2018",
		period:       "12",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2018, 0)

	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestMonthlyInvalidMonthIntXFail(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "44",
		expectedCode: 400,
	}
	tearDown(t)
	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestMonthlyInvalidMonthStringXFail(t *testing.T) {
	var tc = testCase{
		year:         "2018",
		period:       "Q",
		expectedCode: 400,
	}
	tearDown(t)
	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestMonthlyInvalidYearStringXFail(t *testing.T) {
	var tc = testCase{
		year:         "Q",
		period:       "1",
		expectedCode: 400,
	}
	tearDown(t)
	assertMonthlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ42017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q4",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 10, 3, 2017, 4) // Oct-Dec 2017 status 4
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ32017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q3",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 7, 3, 2017, 4)
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ22017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q2",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 4, 3, 2017, 4) // Jan-Jun 2017
	assertQuarterlyStatusCodeEqual(t, &tc)

}

func TestQuarterlyQ12017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q1",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 3, 2017, 4)
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ42017XFail(t *testing.T) {
	// Already exists
	var tc = testCase{
		year:         "2017",
		period:       "Q1",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 3, 2017, 4)
	setupQuarterlyTables(t, 1, 1, 2017, 4)
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ42010XFail(t *testing.T) {
	// No monthly batches exist for Oct-Dec 2010
	var tc = testCase{
		year:         "2019",
		period:       "Q4",
		expectedCode: 400,
	}
	tearDown(t)
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ12015XFail(t *testing.T) {
	// 2 monthly batches exist for Jan-Feb 2015. Required 3 monthly batches to continue
	var tc = testCase{
		year:         "2015",
		period:       "Q1",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 2, 2015, 0) // Jan-Feb 2015
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ12016XFail(t *testing.T) {
	// 2 VALID monthly batches exist for Jan-Feb 2016. Required 3 VALID monthly batches to continue
	var tc = testCase{
		year:         "2016",
		period:       "Q1",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 2, 2016, 4) // Jan-Feb 2016
	setupMonthlyTables(t, 3, 1, 2016, 0) // Mar 2016
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyInvalidQuarterIntXFail(t *testing.T) {
	var tc = testCase{
		year:         "2019",
		period:       "4",
		expectedCode: 400,
	}
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyInvalidQuarterXFail(t *testing.T) {
	var tc = testCase{
		year:         "2019",
		period:       "Q5",
		expectedCode: 400,
	}
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestQuarterlyInvalidYearXFail(t *testing.T) {
	var tc = testCase{
		year:         "El is a superstar",
		period:       "Q4",
		expectedCode: 400,
	}
	assertQuarterlyStatusCodeEqual(t, &tc)
}

func TestAnnual2017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2017, 4)  // Jan-Dec 2017 status 4
	setupQuarterlyTables(t, 1, 4, 2017, 4) // Q1-Q4 2017 status 4
	assertAnnualStatusCodeEqual(t, &tc)
}

func TestAnnual2017AlreadyExistsXFail(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2017, 4)  // Jan-Dec 2017 status 4
	setupQuarterlyTables(t, 1, 4, 2017, 4) // Q1-Q4 2017 status 4
	setupAnnualTables(t, 2017, 4)          // Q1-Q4 2017 status 4
	assertAnnualStatusCodeEqual(t, &tc)
}

func TestAnnualZeroMonhlyBatchesXFail(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	assertAnnualStatusCodeEqual(t, &tc)
}

func TestAnnualMonhlyBatchesXFail(t *testing.T) {
	// 12 monthly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 5, 2014, 0) // Jan-May 2014
	assertAnnualStatusCodeEqual(t, &tc)
}

func TestAnnualValidMonhlyBatchesXFail(t *testing.T) {
	// 12 VALID monthly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 5, 2014, 4) // Jan-May 2014
	assertAnnualStatusCodeEqual(t, &tc)
}

func TestAnnualQuarterlyBatchesXFail(t *testing.T) {
	// 4 quarterly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2014, 4) // Jan-Dec 2014
	assertAnnualStatusCodeEqual(t, &tc)
}

func TestAnnualValidQuarterlyBatchesXFail(t *testing.T) {
	// 4 VALID quarterly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2014, 4)  // Jan-Dec 2014
	setupQuarterlyTables(t, 1, 2, 2014, 4) // Q1-Q2 2014 - complete
	setupQuarterlyTables(t, 3, 2, 2014, 0) // Q3-Q4 2014 - open
	assertAnnualStatusCodeEqual(t, &tc)
}

// TODO: Improve final teardown
func TestFinalTearDown(t *testing.T) {
	tearDown(t)
}

func countRows(t *testing.T, tableName string) int {
	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		t.Fatalf(err.Error())
	}

	count, err := dbase.CountRows(tableName)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return count
}

func assertMonthlyStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/monthly/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year, "month": tc.period})

	BatchHandler{}.CreateMonthlyBatchHandler(w, r)

	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\nERROR: %s", w.Body.String())
	}

	// Database Updated
	if tc.expectedCode == 200 {
		assert.Equal(t, 1, countRows(t, "monthly_batch"))
	}
}

func assertQuarterlyStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/quarterly/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year, "quarter": tc.period})

	BatchHandler{}.CreateQuarterlyBatchHandler(w, r)

	// Status Code
	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\nERROR: %s", w.Body.String())
	}

	// Database Updated
	if tc.expectedCode == 200 {
		assert.Equal(t, 1, countRows(t, "quarterly_batch"))
	}
}

func assertAnnualStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/annual/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year})

	BatchHandler{}.CreateAnnualBatchHandler(w, r)

	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\nERROR: %s", w.Body.String())
	}

	// Database Updated
	if tc.expectedCode == 200 {
		assert.Equal(t, 1, countRows(t, "annual_batch"))
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
	// Insert a load of mock data
	month -= 1
	for c := 1; c <= count; c++ {
		batch := types.MonthlyBatch{
			Year:        year,
			Month:       month + c,
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
	quarter -= 1
	for c := 1; c <= count; c++ {
		batch := types.QuarterlyBatch{
			Quarter:     quarter + c,
			Year:        year,
			Status:      status,
			Description: "Mock data for Testing",
		}
		if err := dbase.CreateQuarterlyBatch(batch); err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func setupAnnualTables(t *testing.T, year, status int) {
	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		t.Fatalf(err.Error())
	}

	annualTable := config.Config.Database.AnnualBatchTable
	if annualTable == "" {
		t.Fatal("annual_batch table configuration not set")
	}

	// Insert a load of mock data and set status to 4
	batch := types.AnnualBatch{
		Year:        year,
		Status:      status,
		Description: "Mock data for Testing",
	}
	if err := dbase.CreateAnnualBatch(batch); err != nil {
		t.Fatalf(err.Error())
	}
}
