package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
)

func init() {
	batchTable = config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		panic("monthly batch table configuration not set")
	}

	quarterlyBatchTable = config.Config.Database.QuarterlyBatchTable
	if quarterlyBatchTable == "" {
		panic("quarterly batch table configuration not set")
	}

	annualBatchTable = config.Config.Database.AnnualBatchTable
	if annualBatchTable == "" {
		panic("annual batch table configuration not set")
	}
}

func (s MySQL) GetAnnualBatches() ([]types.Dashboard, error) {
	// Variables
	var annualBatches []types.Dashboard

	// Get table
	dbAnnual := s.DB.Collection(annualBatchTable)

	// Error handling
	if !dbAnnual.Exists() {
		log.Error().Str("table", annualBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", annualBatchTable)
	}

	// Get values
	res := dbAnnual.Find()
	err := res.All(&annualBatches)

	for i := range annualBatches {
		annualBatches[i].Type = "Annually"
	}

	// Error handling
	if err != nil {
		log.Debug().
			Msg("Get Annual Batches: " + err.Error())
		return nil, err
	}

	return annualBatches, nil
}

func (s MySQL) GetQuarterlyBatches() ([]types.Dashboard, error) {
	// Variables
	var quarterlyBatches []types.Dashboard

	// Get table
	dbQuarter := s.DB.Collection(quarterlyBatchTable)

	// Error handling
	if !dbQuarter.Exists() {
		log.Error().Str("table", quarterlyBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", quarterlyBatchTable)
	}

	// Get values
	res := dbQuarter.Find()
	err := res.All(&quarterlyBatches)

	for i := range quarterlyBatches {
		quarterlyBatches[i].Type = "Quarterly"
	}

	// Error handling
	if err != nil {
		log.Debug().
			Msg("Get Quarterly Batches: " + err.Error())
		return nil, err
	}

	return quarterlyBatches, nil
}

func (s MySQL) GetMonthlyBatches() ([]types.Dashboard, error) {
	// GB and NI Variables
	var monthlyBatches []types.Dashboard

	// Get GB and NI tables
	db := s.DB.Collection(batchTable)

	// Error handling
	if !db.Exists() {
		log.Error().Str("table", batchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", batchTable)
	}

	// Get GB and NI values
	res := db.Find()
	err := res.All(&monthlyBatches)

	for i := range monthlyBatches {
		monthlyBatches[i].Type = "Monthly"
	}

	// Error handling
	if err != nil {
		log.Debug().
			Msg("Get Monthly Batches : " + err.Error())
		return nil, err
	}

	return monthlyBatches, nil
}
