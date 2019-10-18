package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
)

var addressesTable string

func init() {
	addressesTable = config.Config.Database.AddressesTable
	if addressesTable == "" {
		panic("addresses table configuration not set")
	}
}

func (s MySQL) DeleteAddressesData(name string) error {
	col := s.DB.Collection(addressesTable)
	res := col.Find("table_name", name)
	if res == nil {
		return nil
	}
	if err := res.Delete(); err != nil {
		return err
	}
	return nil
}

func (s MySQL) insertAddressesData(tx sqlbuilder.Tx, addresses map[string]string) error {

	col := tx.Collection(addressesTable)
	_, err := col.Insert(addresses)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) PersistAddressDataset(tmpFile string) error {

	startTime := time.Now()
	log.Debug().
		Str("tableName", addressesTable).
		Str("fileName", tmpFile).
		Msg("Starting persistence into DB")

	query := "LOAD DATA INFILE '" + tmpFile + "' REPLACE INTO TABLE " + addressesTable + " FIELDS TERMINATED BY ','"
	res, err := s.DB.Exec(query)
	if err != nil {
		log.Error().
			Err(err).
			Msg("LOAD Data failed")
		return fmt.Errorf("load data statement failed, error: %s", err)
	}

	cnt, _ := res.RowsAffected()

	log.Debug().
		TimeDiff("elapsedTime", time.Now(), startTime).
		Int64("Rows inserted", cnt).
		Msg("Data persisted")

	//tx, err := s.DB.NewTx(nil)
	//if err != nil {
	//	log.Error().
	//		Err(err).
	//		Msg("Start transaction failed")
	//	return fmt.Errorf("cannot start a transaction, error: %s", err)
	//}

	//var f = DBAudit{s}
	//if err := f.AuditFileUploadEvent(tx, d); err != nil {
	//	log.Error().
	//		Err(err).
	//		Msg("AuditFileUpload failed")
	//	return fmt.Errorf("AuditFileUpload, error: %s", err)
	//}
	//
	//if err := tx.Commit(); err != nil {
	//	log.Error().
	//		Err(err).
	//		Msg("Commit transaction failed")
	//	return fmt.Errorf("commit failed, error: %s", err)
	//}

	return nil
}
