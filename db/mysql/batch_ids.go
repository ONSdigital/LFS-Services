package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
	"upper.io/db.v3"
)

type DBId struct {
	MySQL
}

//var surveyAuditTable string

func init() {
	batchTable = config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		panic("monthly batch table configuration not set")
	}
	//
	//gbBatchTable = config.Config.Database.GbBatchTable
	//if gbBatchTable == "" {
	//	panic("gb batch table configuration not set")
	//}
	//
	niBatchTable = config.Config.Database.NiBatchTable
	if niBatchTable == "" {
		panic("ni batch table configuration not set")
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

func (s MySQL) GetIdsByYear(year types.Year) ([]types.YearID, error) {
	// Variables
	var yearIDs []types.YearID

	// Get table
	dbAnnual := s.DB.Collection(annualBatchTable)

	// Error handling
	if !dbAnnual.Exists() {
		log.Error().Str("table", annualBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", annualBatchTable)
	}

	// Get values
	res := dbAnnual.Find(db.Cond{"year": year})
	err := res.All(&yearIDs)

	// Error handling
	if err != nil {
		log.Debug().
			Int("year", int(year)).
			Msg("Get Annual Batch IDs failed: " + err.Error())
		return nil, err
	}

	return yearIDs, nil
}

func (s MySQL) GetIdsByQuarter(year types.Year, quarter types.Quarter) ([]types.QuarterID, error) {
	// Variables
	var quarterID []types.QuarterID

	// Get table
	dbQuarter := s.DB.Collection(quarterlyBatchTable)

	// Error handling
	if !dbQuarter.Exists() {
		log.Error().Str("table", quarterlyBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", quarterlyBatchTable)
	}

	// Get values
	res := dbQuarter.Find(db.Cond{"year": year, "quarter": quarter})
	err := res.All(&quarterID)

	// Error handling
	if err != nil {
		log.Debug().
			Int("year", int(year)).
			Msg("Get Quarterly Batch IDs failed: " + err.Error())
		return nil, err
	}

	return quarterID, nil
}

func (s MySQL) GetIdsByMonth(year types.Year, month types.Month) ([]types.MonthID, error) {
	// Variables
	var monthID []types.MonthID

	// Get table
	dbMonth := s.DB.Collection(batchTable)

	// Error handling
	if !dbMonth.Exists() {
		log.Error().Str("table", batchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", batchTable)
	}

	// Get values
	res := dbMonth.Find(db.Cond{"year": year, "month": month})
	err := res.All(&monthID)

	// Error handling
	if err != nil {
		log.Debug().
			Int("year", int(year)).
			Msg("Get Monthly Batch IDs failed: " + err.Error())
		return nil, err
	}

	return monthID, nil
}

func (s MySQL) GetNIIds(year types.Year, month types.Month) ([]types.NIID, error) {
	// Variables
	var niID []types.NIID

	// Get table
	dbNI := s.DB.Collection(niBatchTable)

	// Error handling
	if !dbNI.Exists() {
		log.Error().Str("table", niBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", niBatchTable)
	}

	// Get values
	res := dbNI.Find(db.Cond{"year": year, "month": month})
	err := res.All(&niID)

	// Error handling
	if err != nil {
		log.Debug().
			Int("year", int(year)).
			Msg("Get Monthly Batch IDs failed: " + err.Error())
		return nil, err
	}

	return niID, nil
}
