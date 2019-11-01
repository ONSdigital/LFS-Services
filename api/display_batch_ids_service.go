package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/types"
)

func (i IdHandler) GetIdsForYear(year types.Year) ([]types.YearID, error) {
	// Database connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// Retrieve table values
	res, err := dbase.GetIdsByYear(year)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i IdHandler) GetIdsForQuarter(year types.Year, quarter types.Quarter) ([]types.QuarterID, error) {
	// Database connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// Retrieve table values
	res, err := dbase.GetIdsByQuarter(year, quarter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i IdHandler) GetIdsForMonth(year types.Year, month types.Month) ([]types.MonthID, error) {
	// Error capture
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("the month value is %d, must be between 1 and 12", month)
	}

	// Database connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// Retrieve table values
	res, err := dbase.GetIdsByMonth(year, month)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i IdHandler) GetIdsForNI(year types.Year, month types.Month) ([]types.NIID, error) {
	// Error capture
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("the month value is %d, must be between 1 and 12", month)
	}

	// Database connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// Retrieve table values
	res, err := dbase.GetNIIds(year, month)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i IdHandler) GetIdsForGB(year types.Year, month types.Month, week types.Week) ([]types.GBID, error) {
	// Error capture
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("the month value is %d, must be between 1 and 12", month)
	}

	if week < 1 || week > 5 {
		return nil, fmt.Errorf("the week value is %d, must be between 1 and 5", week)
	}

	// Database connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// Retrieve table values
	res, err := dbase.GetGBIds(year, month, week)
	if err != nil {
		return nil, err
	}
	return res, nil
}
