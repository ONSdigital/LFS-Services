package postgres

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
	"upper.io/db.v3"
)

var batchTable string
var gbBatchTable string
var niBatchTable string
var quarterlyBatchTable string
var annualBatchTable string

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

func (s Postgres) MonthlyBatchExists(month, year int) bool {
	col := s.DB.Collection(batchTable)
	res := col.Find(db.Cond{"month": month, "year": year})

	type R struct {
		month int
		year  int
	}
	var result R
	if err := res.One(&result); err != nil {
		log.Debug().
			Int("month", month).
			Int("year", year).
			Msg("Batch does not exist")
		return false
	}

	log.Warn().
		Int("month", month).
		Int("year", year).
		Msg("Monthly batch check - Batch already exists")

	return true
}

func (s Postgres) AnnualBatchExists(year int) bool {
	col := s.DB.Collection(annualBatchTable)
	//res := col.Find("year", year)
	res := col.Find(db.Cond{"year": year})

	type R struct {
		year int
	}
	var result R
	if err := res.One(&result); err != nil {
		log.Debug().
			Int("year", year).
			Msg("Batch does not exist")
		return false
	}

	log.Warn().
		Int("year", year).
		Msg("Annual batch check - Batch already exists")

	return true
}

func (s Postgres) QuarterBatchExists(quarter, year int) bool {
	col := s.DB.Collection(quarterlyBatchTable)
	res := col.Find(db.Cond{"quarter": quarter, "year": year})

	type R struct {
		month int
		year  int
	}
	var result R
	if err := res.One(&result); err != nil {
		log.Debug().
			Int("quarter", quarter).
			Int("year", year).
			Msg("Batch does not exist")
		return false
	}

	log.Warn().
		Int("quarter", quarter).
		Int("year", year).
		Msg("Monthly batch check - Batch already exists")

	return true
}

func (s Postgres) ValidateMonthsForQuarterlyBatch(period, year int) ([]types.MonthlyBatch, error) {
	var entireResult []types.MonthlyBatch
	var validResult []types.MonthlyBatch
	var months []int
	switch period {
	case 1:
		months = append(months, 1, 2, 3)
	case 2:
		months = append(months, 4, 5, 6)
	case 3:
		months = append(months, 7, 8, 9)
	case 4:
		months = append(months, 10, 11, 12)
	}

	col := s.DB.Collection(batchTable)
	if !col.Exists() {
		return nil, fmt.Errorf("table: %s does not exist", batchTable)
	}

	// Confirm 3 months exist
	threeRes := col.Find(db.Cond{"year": year, "month": months})
	threeErr := threeRes.All(&entireResult)
	if threeErr != nil {
		return nil, threeErr
	}

	total, err := threeRes.Count()
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return nil, fmt.Errorf("no monthly batches exist for Q%d, %d. Required 3 monthly batches to continue", period, year)
	}

	if total < 3 {
		return entireResult, fmt.Errorf(
			"%d monthly batches exist for Q%d, %d. Required 3 monthly batches to continue",
			total, period, year)
	}

	// Get results for a valid quarter where status is complete
	validRes := col.Find(db.Cond{"year": year, "status": 4, "month": months})
	validErr := validRes.All(&validResult)
	if validErr != nil {
		return nil, validErr
	}

	// Count results
	total, err = validRes.Count()
	if err != nil {
		return nil, err
	}

	if total < 3 {
		// Return valid and invalid results for a the quarter
		return entireResult, fmt.Errorf(
			"%d valid monthly batches exist for Q%d, %d. Required 3 valid monthly batches to continue",
			total, period, year)
	}

	return nil, nil
}

func (s Postgres) ValidateMonthsForAnnualBatch(year int) ([]types.MonthlyBatch, error) {
	var entireResult []types.MonthlyBatch
	var validResult []types.MonthlyBatch

	col := s.DB.Collection(batchTable)
	if !col.Exists() {
		return nil, fmt.Errorf("table: %s does not exist", batchTable)
	}

	// confirm 12 months exist
	res := col.Find(db.Cond{"year": year})
	err := res.All(&entireResult)
	if err != nil {
		return nil, err
	}

	total, countErr := res.Count()
	if countErr != nil {
		return nil, err
	}

	if total == 0 {
		return nil, fmt.Errorf("no monthly batches exist for %d. Required 12 monthly batches to continue", year)
	}

	if total < 12 {
		return entireResult, fmt.Errorf(
			"%d monthly batches exist for %d. Required 12 monthly batches to continue",
			total, year)
	}

	// Get results for a 12 valid months where status is complete
	validRes := col.Find(db.Cond{"year": year, "status": 4})
	validErr := validRes.All(&validResult)
	if validErr != nil {
		return nil, validErr
	}

	// Count results
	validTotal, validErr := validRes.Count()
	if validErr != nil {
		return nil, err
	}

	if validTotal != 12 {
		// Return valid and invalid results for a the quarter
		return entireResult, fmt.Errorf(
			"%d valid monthly batches exist for %d. Required 12 valid monthly batches to continue",
			total, year)
	}

	return nil, nil
}

func (s Postgres) ValidateQuartersForAnnualBatch(year int) ([]types.QuarterlyBatch, error) {
	var entireResult []types.QuarterlyBatch
	var validResult []types.QuarterlyBatch

	col := s.DB.Collection(quarterlyBatchTable)
	if !col.Exists() {
		return nil, fmt.Errorf("table: %s does not exist", quarterlyBatchTable)
	}

	// Confirm 4 quarters exist
	entireRes := col.Find(db.Cond{"year": year})
	entireErr := entireRes.All(&entireResult)
	if entireErr != nil {
		return nil, entireErr
	}

	total, err := entireRes.Count()
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return nil, fmt.Errorf("no quarterly batches exist for %d. Required 4 quarterly batches to continue", year)
	}

	if total < 4 {
		return entireResult, fmt.Errorf(
			"%d quarterly batches exist for %d. Required 4 quarterly batches to continue",
			total, year)
	}

	// Get results for a valid quarter where status is complete
	validRes := col.Find(db.Cond{"year": year, "status": 4})
	validErr := validRes.All(&validResult)
	if validErr != nil {
		return nil, validErr
	}

	// Count results
	total, err = validRes.Count()
	if err != nil {
		return nil, err
	}

	if total != 4 {
		// Return valid and invalid results for a the quarter
		return entireResult, fmt.Errorf(
			"%d valid quarterly batches exist for %d. Required 4 valid quarterly batches to continue",
			total, year)
	}

	return nil, nil
}

func (s Postgres) CreateMonthlyBatch(batch types.MonthlyBatch) error {
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
		log.Error().
			Err(err).
			Msg("Cannot insert into " + batchTable)
		return fmt.Errorf("insert into %s failed, error: %s", batchTable, err)
	}

	niBatch := tx.Collection(niBatchTable)

	var ni types.NIBatchItem
	ni.Month = batch.Month
	ni.Year = batch.Year
	ni.Status = batch.Status
	ni.Id = int(batchId.(int64))
	_, err = niBatch.Insert(ni)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot insert into " + niBatchTable)
		return fmt.Errorf("insert into %s failed, error: %s", niBatchTable, err)
	}

	cnt := 4
	if batch.Month%3 == 0 {
		cnt = 5
	}

	// get week number - if % 3 then 5 weeks else 4
	weekNo := 1
	for i := 1; i < batch.Month; i++ {
		if i%3 == 0 {
			weekNo = weekNo + 5
		} else {
			weekNo = weekNo + 4
		}
	}

	gbBatch := tx.Collection(gbBatchTable)

	for i := 0; i < cnt; i++ {
		var gb types.GBBatchItem

		gb.Month = batch.Month
		gb.Year = batch.Year
		gb.Status = batch.Status
		gb.Week = weekNo
		gb.Id = int(batchId.(int64))
		_, err = gbBatch.Insert(gb)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Cannot insert into " + gbBatchTable)
			return fmt.Errorf("insert into %s failed, error: %s", gbBatchTable, err)
		}
		weekNo++
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	return nil
}

func (s Postgres) CreateQuarterlyBatch(batch types.QuarterlyBatch) error {
	// Create new transaction
	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	// Insert into quarterly_batch
	b := tx.Collection(quarterlyBatchTable)
	_, err = b.Insert(batch)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot insert into " + quarterlyBatchTable)
		return fmt.Errorf("insert into %s failed, error: %s", quarterlyBatchTable, err)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	return nil
}

func (s Postgres) CreateAnnualBatch(batch types.AnnualBatch) error {
	// Create new transaction
	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	// Insert into annual_batch
	b := tx.Collection(annualBatchTable)
	_, err = b.Insert(batch)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot insert into " + annualBatchTable)
		return fmt.Errorf("insert into %s failed, error: %s", annualBatchTable, err)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	return nil
}

func (s Postgres) FindGBBatchInfo(week, year int) (types.GBBatchItem, error) {
	batchInfo := s.DB.Collection(gbBatchTable)
	//res := batchInfo.Find("week", week, "year", year)
	res := batchInfo.Find(db.Cond{"week": week, "year": year})

	var result types.GBBatchItem
	if err := res.One(&result); err != nil {
		log.Debug().
			Int("week", week).
			Int("year", year).
			Msg("GB batch does not exist")
		return types.GBBatchItem{}, err
	}
	return result, nil

}

func (s Postgres) FindNIBatchInfo(month, year int) (types.NIBatchItem, error) {
	batchInfo := s.DB.Collection(niBatchTable)
	//res := batchInfo.Find("month", month, "year", year)
	res := batchInfo.Find(db.Cond{"month": month, "year": year})

	var result types.NIBatchItem
	if err := res.One(&result); err != nil {
		log.Debug().
			Int("month", month).
			Int("year", year).
			Msg("Batch does not exist")
		return types.NIBatchItem{}, err
	}
	return result, nil
}

func (s Postgres) UpdateNIMonthlyStatus(week, month, status int) error {
	batchInfo := s.DB.Collection(niBatchTable)
	//res := batchInfo.Find("week", week, "mont", month)
	res := batchInfo.Find(db.Cond{"week": week, "month": month})

	var result types.NIBatchItem
	if err := res.One(&result); err != nil {
		log.Debug().
			Int("week", week).
			Int("month", month).
			Msg("NI batch does not exist")
		return err
	}
	result.Status = status
	if err := res.Update(result); err != nil {
		log.Debug().
			Int("week", week).
			Int("month", month).
			Msg("NI batch update failed")
		return err
	}
	return nil
}

func (s Postgres) updateGBBatch(week, year, status int) error {
	batchInfo := s.DB.Collection(gbBatchTable)
	//res := batchInfo.Find("week", week, "year", year)
	res := batchInfo.Find(db.Cond{"week": week, "year": year})

	var result types.GBBatchItem

	if err := res.One(&result); err != nil {
		log.Debug().Int("week", week).Int("year", year).Msg("Batch does not exist")
		return err
	}

	result.Status = status

	if err := res.Update(result); err != nil {
		log.Debug().Int("week", week).Int("year", year).Msg("GB batch update failed")
		return err
	}
	return nil
}

func (s Postgres) updateNIBatch(month, year, status int) error {
	batchInfo := s.DB.Collection(niBatchTable)
	//res := batchInfo.Find("month", month, "year", year)
	res := batchInfo.Find(db.Cond{"month": month, "year": year})
	var result types.NIBatchItem

	if err := res.One(&result); err != nil {
		log.Debug().Int("month", month).Int("year", year).Msg("Batch does not exist")
		return err
	}

	result.Status = status

	if err := res.Update(result); err != nil {
		log.Debug().Int("week", month).Int("year", year).Msg("GB batch update failed")
		return err
	}
	return nil
}
