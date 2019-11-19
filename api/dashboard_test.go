package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"services/config"
	"testing"
)

var batchTable string
var quarterlyBatchTable string
var annualBatchTable string

func init() {
	batchTable = config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		panic("monthly batch table configuration not set")
	}

	quarterlyBatchTable = config.Config.Database.QuarterlyBatchTable
	if quarterlyBatchTable == "" {
		panic("quarterly batch table configuration not set")
	}

	annualBatchTable = config.Config.Database.AnnualBatchTable
	if annualBatchTable == "" {
		panic("annual batch table configuration not set")
	}
}

func TestDashboardSuccess1(t *testing.T) {
	var tc = testCase{
		expectedCode: 200,
	}

	tearDown(t)

	// Set up random amount of batches
	setupMonthlyTables(t, 1, 12, 2014, 4)
	setupMonthlyTables(t, 1, 7, 2015, 4)
	setupMonthlyTables(t, 8, 5, 2015, 0)
	setupMonthlyTables(t, 1, 1, 2016, 0)
	setupQuarterlyTables(t, 1, 4, 2014, 4)
	setupQuarterlyTables(t, 1, 2, 2015, 4)
	setupQuarterlyTables(t, 3, 2, 2015, 0)
	setupAnnualTables(t, 2014, 4)

	monthlyTotal := countRows(t, batchTable)
	quarterlyTotal := countRows(t, quarterlyBatchTable)
	annualTotal := countRows(t, annualBatchTable)

	expectedTotal := monthlyTotal + quarterlyTotal + annualTotal
	tc.expectedItems = expectedTotal

	assertDashboardTest(t, &tc)
}

func TestDashboardSuccess2(t *testing.T) {
	var tc = testCase{
		expectedCode: 200,
	}

	tearDown(t)

	// Set up random amount of batches
	setupMonthlyTables(t, 1, 12, 2014, 4)
	setupQuarterlyTables(t, 1, 4, 2014, 4)
	setupQuarterlyTables(t, 1, 2, 2015, 4)

	monthlyTotal := countRows(t, batchTable)
	quarterlyTotal := countRows(t, quarterlyBatchTable)
	annualTotal := countRows(t, annualBatchTable)

	expectedTotal := monthlyTotal + quarterlyTotal + annualTotal
	tc.expectedItems = expectedTotal

	assertDashboardTest(t, &tc)
}

func TestDashboardSuccess3(t *testing.T) {
	var tc = testCase{
		expectedCode: 200,
	}

	tearDown(t)

	// Set up random amount of batches
	setupMonthlyTables(t, 1, 12, 2014, 4)
	setupQuarterlyTables(t, 1, 4, 2014, 4)
	setupAnnualTables(t, 2014, 4)

	monthlyTotal := countRows(t, batchTable)
	quarterlyTotal := countRows(t, quarterlyBatchTable)
	annualTotal := countRows(t, annualBatchTable)

	expectedTotal := monthlyTotal + quarterlyTotal + annualTotal
	tc.expectedItems = expectedTotal

	assertDashboardTest(t, &tc)
}

func TestDashboardSuccess4(t *testing.T) {
	var tc = testCase{
		expectedCode: 200,
	}

	tearDown(t)

	// Set up random amount of batches
	setupMonthlyTables(t, 1, 3, 2014, 0)

	monthlyTotal := countRows(t, batchTable)
	quarterlyTotal := countRows(t, quarterlyBatchTable)
	annualTotal := countRows(t, annualBatchTable)

	expectedTotal := monthlyTotal + quarterlyTotal + annualTotal
	tc.expectedItems = expectedTotal

	assertDashboardTest(t, &tc)
}

func TestDashboardFinalTearDown(t *testing.T) {
	tearDown(t)
}

func assertDashboardTest(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	DashboardHandler{}.HandleDashboardRequest(w, r)

	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}
	t.Log(">>>>> PASS: Response code from HandleDashboardRequest() received as expected")

	// Assert number of batches returned in JSON
	actual := bytes.Count(w.Body.Bytes(), []byte("{"))
	if !assert.Equal(t, actual, tc.expectedItems) {
		t.Fatalf("\n\n>>>>> Expected %v batches. Actual number of retrieved batches were %v",
			tc.expectedItems,
			actual)
	}
	t.Log(">>>>> PASS: Response returned correct number of batches")
	t.Log("\n")

}
