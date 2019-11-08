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

func (b BatchHandler) generateQuarterBatchId(quarter int, year int, description string) error {
	batch := types.QuarterlyBatch{
		Id:          0,
		Quarter:     quarter,
		Year:        year,
		Status:      0,
		Description: description,
	}

	// Validate quarter
	if quarter < 1 || quarter > 4 {
		return fmt.Errorf("the quarter value is %d, must be between 1 and 4", quarter)
	}

	// Establish DB connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	// Check if quarter batch already exists
	if found := dbase.QuarterBatchExists(quarter, year); found {
		return fmt.Errorf("q%d batch for year %d already exists", quarter, year)
	}

	// Ensure successful monthly exist
	if found := dbase.ValidateMonthsForQuarterlyBatch(quarter, year); !found {
		return fmt.Errorf("3 valid months for year %d required", year)
	}

	// Create shizznizz
	if err = dbase.CreateQuarterlyBatch(batch); err != nil {
		return err
	}

	return nil
}

func (b BatchHandler) generateYearBatchId(year int, description string) error {
	// Set batch variables
	batch := types.AnnualBatch{
		Id:          1,
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
	if found := dbase.ValidateMonthsForAnnualBatch(year); !found {
		return fmt.Errorf("12 valid months for year %d required", year)
	}

	// Ensure successful quarterly exist
	if found := dbase.ValidateQuartersForAnnualBatch(year); !found {
		return fmt.Errorf("4 valid quarters for year %d required", year)
	}

	// Create shizznizz
	if err = dbase.CreateAnnualBatch(batch); err != nil {
		return err
	}

	return nil
}
