package mysql

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
var quartleryBatchTable string
var annualBatchTable string

func init() {
	batchTable = config.Config.Database.MonthlyBatchTable
	if batchTable == "" {
		panic("monthly batch table configuration not set")
	}
	gbBatchTable = config.Config.Database.GbBatchTable
	niBatchTable = config.Config.Database.NiBatchTable
	quartleryBatchTable = config.Config.Database.QuarterlyBatchTable
	annualBatchTable = config.Config.Database.AnnualBatchTable
}

func (s MySQL) FindNIBatchInfo(month, year int) (types.NIBatchItem, error) {
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

func (s MySQL) FindGBBatchInfo(week, year int) (types.GBBatchItem, error) {
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

func (s MySQL) UpdateNIMonthlyStatus(week, month, status int) error {
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
			Int("year", month).
			Msg("NI batch update failed")
		return err
	}
	return nil
}

func (s MySQL) updateGBBatch(week, year, status int) error {
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

func (s MySQL) updateNIBatch(month, year, status int) error {
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

func (s MySQL) MonthlyBatchExists(month, year int) bool {
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

func (s MySQL) ValidateMonthsForAnnualBatch(year int) bool {
	col := s.DB.Collection(batchTable)
	res := col.Find(db.Cond{"year": year, "status": 4})

	total, err := res.Count()

	if err != nil {
		log.Debug().Msg("Fatal: " + err.Error())
	}

	if total != 12 {
		log.Warn().
			Int("year", year).
			Msg("Cannot continue without 12 valid months")
		return false
	}

	log.Debug().
		Int("year", year).
		Msg("Annual batch check - All 12 valid months exist")

	return true
}

func (s MySQL) ValidateQuartersForAnnualBatch(year int) bool {
	col := s.DB.Collection(quartleryBatchTable)
	res := col.Find(db.Cond{"year": year, "status": 4})

	total, err := res.Count()

	if err != nil {
		log.Debug().Msg("Fatal: " + err.Error())
	}

	if total != 4 {
		log.Warn().
			Int("year", year).
			Msg("Cannot continue without 4 valid quarters")
		return false
	}

	log.Debug().
		Int("year", year).
		Msg("Annual batch check - All 4 valid quarters exist")

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
	weekNo := 0
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

func (s MySQL) AnnualBatchExists(year int) bool {
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

func (s MySQL) CreateAnnualBatch(batch types.AnnualBatch) error {
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
			Msg("Cannot insert into " + batchTable)
		return fmt.Errorf("insert into %s failed, error: %s", batchTable, err)
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

func (s MySQL) QuarterBatchExists(quarter, year int) bool {
	col := s.DB.Collection(quartleryBatchTable)
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

func (s MySQL) ValidateMonthsForQuarterlyBatch(period, year int) bool {
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
	res := col.Find(db.Cond{"year": year, "status": 0, "month": months})

	total, err := res.Count()

	if err != nil {
		log.Debug().Msg("Fatal: " + err.Error())
	}

	if total != 3 {
		log.Warn().
			Int("period", period).
			Msg("Cannot continue without 3 valid months")
		return false
	}

	log.Debug().
		Int("year", year).
		Msg("Qarterly batch check - All 3 valid months exist")

	return true
}

func (s MySQL) CreateQuarterlyBatch(batch types.QuarterlyBatch) error {
	// Create new transaction
	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	// Insert into quarterly_batch
	b := tx.Collection(quartleryBatchTable)
	_, err = b.Insert(batch)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot insert into " + batchTable)
		return fmt.Errorf("insert into %s failed, error: %s", batchTable, err)
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
