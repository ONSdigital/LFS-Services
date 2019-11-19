package api

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestMonthlyMay2014Success(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "05",
		expectedCode: 200,
	}
	tearDown(t)
	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestMonthlyJan2016Success(t *testing.T) {
	var tc = testCase{
		year:         "2016",
		period:       "1",
		expectedCode: 200,
	}
	tearDown(t)
	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestMonthlyDec2018Success(t *testing.T) {
	var tc = testCase{
		year:         "2018",
		period:       "12",
		expectedCode: 200,
	}
	tearDown(t)
	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestMonthlyAlreadyExistsXFail(t *testing.T) {
	var tc = testCase{
		year:         "2018",
		period:       "12",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2018, 0)

	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestMonthlyInvalidMonthIntXFail(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "44",
		expectedCode: 400,
	}
	tearDown(t)
	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestMonthlyInvalidMonthStringXFail(t *testing.T) {
	var tc = testCase{
		year:         "2018",
		period:       "Q",
		expectedCode: 400,
	}
	tearDown(t)
	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestMonthlyInvalidYearStringXFail(t *testing.T) {
	var tc = testCase{
		year:         "Q",
		period:       "1",
		expectedCode: 400,
	}
	tearDown(t)
	assertCreateMonthlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ42017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q4",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 10, 3, 2017, 4) // Oct-Dec 2017 status 4
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ32017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q3",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 7, 3, 2017, 4)
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ22017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q2",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 4, 3, 2017, 4) // Jan-Jun 2017
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)

}

func TestQuarterlyQ12017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		period:       "Q1",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 3, 2017, 4)
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
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
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyQ42010XFail(t *testing.T) {
	// No monthly batches exist for Oct-Dec 2010
	var tc = testCase{
		year:         "2019",
		period:       "Q4",
		expectedCode: 400,
	}
	tearDown(t)
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
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
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
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
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyInvalidQuarterIntXFail(t *testing.T) {
	var tc = testCase{
		year:         "2019",
		period:       "4",
		expectedCode: 400,
	}
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyInvalidQuarterXFail(t *testing.T) {
	var tc = testCase{
		year:         "2019",
		period:       "Q5",
		expectedCode: 400,
	}
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestQuarterlyInvalidYearXFail(t *testing.T) {
	var tc = testCase{
		year:         "El is a superstar",
		period:       "Q4",
		expectedCode: 400,
	}
	assertCreateQuarterlyBatchStatusCodeEqual(t, &tc)
}

func TestAnnual2017Success(t *testing.T) {
	var tc = testCase{
		year:         "2017",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2017, 4)  // Jan-Dec 2017 status 4
	setupQuarterlyTables(t, 1, 4, 2017, 4) // Q1-Q4 2017 status 4
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
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
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
}

func TestAnnualZeroMonthlyBatchesXFail(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
}

func TestAnnualMonthlyBatchesXFail(t *testing.T) {
	// 12 monthly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 5, 2014, 0) // Jan-May 2014
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
}

func TestAnnualValidMonthlyBatchesXFail(t *testing.T) {
	// 12 VALID monthly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 5, 2014, 4) // Jan-May 2014
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
}

func TestAnnualQuarterlyBatchesXFail(t *testing.T) {
	// 4 quarterly batches required
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2014, 4) // Jan-Dec 2014
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
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
	assertCreateAnnualBatchStatusCodeEqual(t, &tc)
}

func TestCreateBatchFinalTearDown(t *testing.T) {
	tearDown(t)
}

func assertCreateMonthlyBatchStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/monthly/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year, "month": tc.period})

	BatchHandler{}.CreateMonthlyBatchHandler(w, r)

	// Status code
	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}
	t.Log(">>>>> PASS: Response code from CreateMonthlyBatchHandler() received as expected")

	// Database Updated
	actual := countRows(t, "monthly_batch")
	if tc.expectedCode == 200 {
		if !assert.Equal(t, 1, actual) {
			t.Fatalf("\n\n>>>>> Expected 1 row in monthly_batch. Actual number of row/s updated were %v",
				actual)
		}
		t.Log(">>>>> PASS: monthly_batch table was updated as expected")
	}

	t.Log("\n")
}

func assertCreateQuarterlyBatchStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/quarterly/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year, "quarter": tc.period})

	BatchHandler{}.CreateQuarterlyBatchHandler(w, r)

	// Status Code
	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}
	t.Log(">>>>> PASS: Response code from CreateQuarterlyBatchHandler() received as expected")

	// Database Updated
	actual := countRows(t, "quarterly_batch")
	if tc.expectedCode == 200 {
		if !assert.Equal(t, 1, actual) {
			t.Fatalf("\n\n>>>>> Expected 1 row in quarterly_batch. Actual number of row/s updated were %v", actual)
		}
		t.Log(">>>>> PASS: quarterly_batch table was updated as expected")
	}

	t.Log("\n")
}

func assertCreateAnnualBatchStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/annual/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year})

	BatchHandler{}.CreateAnnualBatchHandler(w, r)

	// Status code
	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}
	t.Log(">>>>> PASS: Response code from CreateAnnualBatchHandler() received as expected")

	// Database Updated
	actual := countRows(t, "annual_batch")
	if tc.expectedCode == 200 {
		if !assert.Equal(t, 1, actual) {
			t.Fatalf("\n\n>>>>> Expected 1 row in annual_batch. Actual number of row/s updated were %v", actual)
		}
		t.Log(">>>>> PASS: annual_batch table was updated as expected")
	}

	t.Log("\n")
}
