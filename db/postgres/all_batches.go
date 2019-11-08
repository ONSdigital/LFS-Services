package postgres

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
	"strconv"
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

func (s Postgres) GetAnnualBatches() ([]types.Dashboard, error) {
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
			Msg("Get Annual Batches error: " + err.Error())
		return nil, err
	}

	return annualBatches, nil
}

func (s Postgres) GetQuarterlyBatches() ([]types.Dashboard, error) {
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

	for i, v := range quarterlyBatches {
		quarterlyBatches[i].Type = "Quarterly"
		quarterlyBatches[i].Period = "Q" + v.Period
	}

	// Error handling
	if err != nil {
		log.Debug().
			Msg("Get Quarterly Batches error: " + err.Error())
		return nil, err
	}

	return quarterlyBatches, nil
}

func (s Postgres) GetMonthlyBatches() ([]types.Dashboard, error) {
	// Variables
	var monthlyBatches []types.Dashboard
	type period struct {
		Month int `db:"month" json:"month"`
	}
	var periods []period

	// Get monthly_batch table
	db := s.DB.Collection(batchTable)

	// Error handling
	if !db.Exists() {
		log.Error().Str("table", batchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", batchTable)
	}

	// Get values and assign
	res := db.Find()
	mthlyBatchesErr := res.All(&monthlyBatches)
	periodsErr := res.All(&periods)

	for i, v := range periods[:] {
		monthlyBatches[i].Type = "Monthly"
		monthlyBatches[i].Period = strconv.Itoa(v.Month)
	}

	// Error handling
	if mthlyBatchesErr != nil {
		log.Debug().
			Msg("Get Monthly Batches error: " + mthlyBatchesErr.Error())
		return nil, mthlyBatchesErr
	}
	if periodsErr != nil {
		log.Debug().
			Msg("Get periods error: " + periodsErr.Error())
		return nil, periodsErr
	}

	return monthlyBatches, nil
}
