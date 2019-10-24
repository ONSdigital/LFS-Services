package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
)

var batchTable string
var gbBatchTable string
var niBatchTable string

func init() {
	batchTable = config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		panic("monthly batch table configuration not set")
	}
	gbBatchTable = config.Config.Database.GbBatchTable
	niBatchTable = config.Config.Database.NiBatchTable
}

func (s MySQL) MonthlyBatchExists(month, year int) bool {
	col := s.DB.Collection(batchTable)
	res := col.Find("month", month, "year", year)

	type R struct {
		month int
		year  int
	}
	var result R
	err := res.One(&result)

	if err != nil {
		return false
	}

	if res == nil {
		return false
	}
	return true
}

func (s MySQL) CreateMonthlyBatch(batch types.MonthlyBatch) error {

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	b := tx.Collection(batchTable)
	batchId, err := b.Insert(batch)
	if err != nil {
		return err
	}

	niBatch := tx.Collection(niBatchTable)

	var ni types.NIBatchItem
	ni.Month = batch.Month
	ni.Year = batch.Year
	ni.Status = batch.Status
	ni.Id = int(batchId.(int64))
	_, err = niBatch.Insert(ni)
	if err != nil {
		return err
	}

	cnt := 4
	if batch.Month%3 == 0 {
		cnt = 5
	}
	gbBatch := tx.Collection(gbBatchTable)

	for i := 0; i < cnt; i++ {
		var gb types.GBBatchItem
		gb.Month = batch.Month
		gb.Year = batch.Year
		gb.Status = batch.Status
		gb.Week = i + 1
		gb.Id = int(batchId.(int64))
		_, err = gbBatch.Insert(gb)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	return nil
}
