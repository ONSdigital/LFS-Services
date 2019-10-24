package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/types"
)

func (h RestHandlers) generateMonthBatchId(month int, year int, description string) error {
	// Call batch service to validate
	// Get user creds from database
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	found := dbase.MonthlyBatchExists(month, year)
	if found {
		return fmt.Errorf("monthly batch for month %d, year %d already exists", month, year)
	}

	batch := types.MonthlyBatch{
		Id:          0,
		Year:        year,
		Month:       month,
		Status:      0,
		Description: description,
	}

	err = dbase.CreateMonthlyBatch(batch)
	if err != nil {
		return err
	}

	return nil
}

func (h RestHandlers) generateQuarterBatchId(quarter string, year int) error {
	// Call batch service to validate

	return nil
}

func (h RestHandlers) generateYearBatchId(year int) error {
	// Call batch service to validate

	return nil
}
