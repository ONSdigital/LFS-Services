package api

import (
	"github.com/rs/zerolog/log"
	"services/db"
	"services/types"
)

func (d DashboardHandler) GetDashboardInfo() ([]types.Dashboard, error) {
	// Database connection
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
	}

	// Retrieve annual table values
	annualRes, err := dbase.GetAnnualBatches()
	if err != nil {
		return nil, err
	}

	// Retrieve quarterly table values
	qtrRes, err := dbase.GetQuarterlyBatches()
	if err != nil {
		return nil, err
	}

	// Retrieve monthly table values
	mthRes, err := dbase.GetMonthlyBatches()
	if err != nil {
		return nil, err
	}

	// Combine results and return
	for _, a := range annualRes {
		qtrRes = append(qtrRes, a)
	}

	for _, a := range qtrRes {
		mthRes = append(mthRes, a)
	}

	return mthRes, nil
}
