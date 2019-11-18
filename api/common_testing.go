package api

import (
	"services/config"
	"services/db"
	"services/types"
	"testing"
)

type testCase struct {
	year          string
	period        string
	expectedCode  int
	expectedItems int
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
