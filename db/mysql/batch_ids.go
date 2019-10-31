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
	//batchTable = config.Config.Database.MonthlyBatchTable
	//if batchTable == "" {
	//	panic("monthly batch table configuration not set")
	//}
	//
	//gbBatchTable = config.Config.Database.GbBatchTable
	//if gbBatchTable == "" {
	//	panic("gb batch table configuration not set")
	//}
	//
	//niBatchTable = config.Config.Database.NiBatchTable
	//if niBatchTable == "" {
	//	panic("ni batch table configuration not set")
	//}
	//
	//quarterlyBatchTable = config.Config.Database.QuarterlyBatchTable
	//if quarterlyBatchTable == "" {
	//	panic("quarterly batch table configuration not set")
	//}

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
