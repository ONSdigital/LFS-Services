package api

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestMonthBatchInfoMaySuccess(t *testing.T) {
	var tc = testCase{
		year:          "2014",
		period:        "05",
		expectedCode:  200,
		expectedItems: 6, // 5 weekly GB and 1 NI
	}

	tearDown(t)
	setupMonthlyTables(t, 5, 1, 2014, 4)
	assertMonthlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestMonthBatchInfoJanuarySuccess(t *testing.T) {
	var tc = testCase{
		year:          "2014",
		period:        "1",
		expectedCode:  200,
		expectedItems: 5, // 4 weekly GB and 1 NI
	}

	tearDown(t)
	setupMonthlyTables(t, 1, 1, 2014, 4)
	assertMonthlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestMonthBatchInfoFebruarySuccess(t *testing.T) {
	var tc = testCase{
		year:          "2014",
		period:        "2",
		expectedCode:  200,
		expectedItems: 6, // 5 weekly GB and 1 NI
	}

	tearDown(t)
	setupMonthlyTables(t, 2, 1, 2014, 4)
	assertMonthlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestMonthBatchInfoFail(t *testing.T) {
	var tc = testCase{
		year:         "2015",
		period:       "05",
		expectedCode: 400,
	}
	tearDown(t)
	assertMonthlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestMonthBatchInfoFailMonth(t *testing.T) {
	var tc = testCase{
		year:         "2015",
		period:       "13",
		expectedCode: 400,
	}
	tearDown(t)
	assertMonthlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestMonthBatchInfoFailYear(t *testing.T) {
	var tc = testCase{
		year:         "-1",
		period:       "2",
		expectedCode: 400,
	}
	tearDown(t)
	assertMonthlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestQuarterBatchInfoSuccess(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "Q1",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 3, 2014, 4)
	setupQuarterlyTables(t, 1, 1, 2014, 4)
	assertQuarterlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestQuarterBatchInfoFail(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "Q1",
		expectedCode: 400,
	}
	tearDown(t)
	assertQuarterlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestQuarterBatchInfoFailPeriod(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "Q5",
		expectedCode: 400,
	}
	tearDown(t)
	assertQuarterlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestQuarterBatchInfoFailPeriodType(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		period:       "Q",
		expectedCode: 400,
	}
	tearDown(t)
	assertQuarterlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestQuarterBatchInfoFailYear(t *testing.T) {
	var tc = testCase{
		year:         "L",
		period:       "Q3",
		expectedCode: 400,
	}
	tearDown(t)
	assertQuarterlyBatchInfoStatusCodeEqual(t, &tc)
}

func TestAnnualBatchInfoSuccess(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		expectedCode: 200,
	}
	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2014, 4)
	setupQuarterlyTables(t, 1, 4, 2014, 4)
	setupAnnualTables(t, 2014, 4)
	assertAnnualBatchInfoStatusCodeEqual(t, &tc)
}

func TestAnnualBatchInfoFail(t *testing.T) {
	var tc = testCase{
		year:         "2014",
		expectedCode: 400,
	}
	tearDown(t)
	assertAnnualBatchInfoStatusCodeEqual(t, &tc)
}

func TestAnnualBatchInfoFailYear(t *testing.T) {
	// TODO: This isn't failing as expected >:-|
	var tc = testCase{
		year:         "O",
		expectedCode: 400,
	}
	tearDown(t)
	assertAnnualBatchInfoStatusCodeEqual(t, &tc)
}

func assertMonthlyBatchInfoStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/display/monthly/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year, "month": tc.period})

	IdHandler{}.HandleMonthlyBatchIdsRequest(w, r)

	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}

	// Assert number of batches returned in JSON
	if w.Code == 200 {
		if !assert.Equal(t, bytes.Count(w.Body.Bytes(), []byte("{")), tc.expectedItems) {
			t.Fatalf("\n\n>>>>> Expected %v batches. Actual number of retrieved batches were %v",
				tc.expectedItems,
				bytes.Count(w.Body.Bytes(), []byte("{")))
		}
	}

	t.Log("\n")

}

func assertQuarterlyBatchInfoStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/display/quarterly/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year, "quarter": tc.period})

	IdHandler{}.HandleQuarterlyBatchIdsRequest(w, r)

	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}

	t.Log("\n")

}

func assertAnnualBatchInfoStatusCodeEqual(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("POST", "/batches/display/annual/", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"year": tc.year})

	IdHandler{}.HandleAnnualBatchIdsRequest(w, r)

	if !assert.Equal(t, tc.expectedCode, w.Code) {
		t.Fatalf("\n\n>>>>>ERROR: %s", w.Body.String())
	}

	t.Log("\n")
}
