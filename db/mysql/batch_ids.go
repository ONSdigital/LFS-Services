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

	gbBatchTable = config.Config.Database.GbBatchTable
	if gbBatchTable == "" {
		panic("gb batch table configuration not set")
	}

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
	// GB and NI Variables
	var gbMonthID []types.MonthID
	var niMonthID []types.MonthID

	// Get GB and NI tables
	dbGB := s.DB.Collection(gbBatchTable)
	dbNI := s.DB.Collection(niBatchTable)

	// Error handling
	if !dbGB.Exists() {
		log.Error().Str("table", gbBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", gbBatchTable)
	}
	if !dbNI.Exists() {
		log.Error().Str("table", niBatchTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", niBatchTable)
	}

	// Get GB and NI values
	gbRes := dbGB.Find(db.Cond{"year": year, "month": month}).OrderBy("week")
	gbErr := gbRes.All(&gbMonthID)

	niRes := dbNI.Find(db.Cond{"year": year, "month": month})
	niErr := niRes.All(&niMonthID)

	// Populate Location type
	for i := range gbMonthID {
		gbMonthID[i].Loc = "GB"
	}
	for i := range niMonthID {
		niMonthID[i].Loc = "NI"
	}

	// Error handling
	if gbErr != nil {
		log.Debug().
			Int("year", int(year)).
			Msg("Get Monthly Batch IDs failed: " + gbErr.Error())
		return nil, gbErr
	}
	if niErr != nil {
		log.Debug().
			Int("year", int(year)).
			Msg("Get Monthly Batch IDs failed: " + niErr.Error())
		return nil, niErr
	}

	// Combine and return results
	for _, a := range niMonthID {
		gbMonthID = append(gbMonthID, a)
	}

	return gbMonthID, nil
}
