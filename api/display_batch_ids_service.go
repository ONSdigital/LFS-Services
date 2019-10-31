package api

import (
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