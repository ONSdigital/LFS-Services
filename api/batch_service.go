package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/types"
)

func (b BatchHandler) generateMonthBatchId(month int, year int, description string) error {

	if month < 1 || month > 12 {
		return fmt.Errorf("the month value is %d, must be between 1 and 12", month)
	}

	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if found := dbase.MonthlyBatchExists(month, year); found {
		return fmt.Errorf("monthly batch for month %d, year %d already exists", month, year)
	}

	batch := types.MonthlyBatch{
		//Id:          0,
		Year:        year,
		Month:       month,
		Status:      0,
		Description: description,
	}

	if err = dbase.CreateMonthlyBatch(batch); err != nil {
		return err
	}

	return nil
}

func (b BatchHandler) generateQuarterBatchId(quarter string, year int) error {
	// Call batch service to validate

	return nil
}

func (b BatchHandler) generateYearBatchId(year int, description string) error {
	// Set batch variable
	batch := types.AnnualBatch{
		Id:          0,
		Year:        year,
		Status:      0,
		Description: description,
	}

	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	// Check if year batch already exists
	if found := dbase.AnnualBatchExists(year); found {
		return fmt.Errorf("annual batch for year %d already exists", year)
	}

	// Ensure successful monthly exist
	if found := dbase.SuccessfulMonthlyBatchesExist(year); !found {
		return fmt.Errorf("12 valid months for year %d required", year)
	}

	// Ensure successful quarterly exist
	if found := dbase.SuccessfulQuarterlyBatchesExist(year); !found {
		return fmt.Errorf("4 valid quarters for year %d required", year)
	}

	// Create shizznizz
	if err = dbase.CreateAnnualBatch(batch); err != nil {
		return err
	}

	return nil
}
