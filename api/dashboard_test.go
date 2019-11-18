package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestDashboardSuccess1(t *testing.T) {
	var tc = testCase{
		expectedCode:  200,
		expectedItems: 34, // TODO: Count db and assign here/there
	}

	tearDown(t)
	setupMonthlyTables(t, 1, 12, 2014, 4)
	setupMonthlyTables(t, 1, 7, 2015, 4)
	setupMonthlyTables(t, 8, 5, 2015, 0)
	setupMonthlyTables(t, 1, 1, 2016, 0)
	setupQuarterlyTables(t, 1, 4, 2014, 4)
	setupQuarterlyTables(t, 1, 2, 2015, 4)
	setupQuarterlyTables(t, 3, 2, 2015, 0)
	setupAnnualTables(t, 2014, 4)
	assertDashboardTest(t, &tc)
}

func assertDashboardTest(t *testing.T, tc *testCase) {
	r := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	DashboardHandler{}.HandleDashboardRequest(w, r)

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
