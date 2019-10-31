package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/importdata/csv"
	"services/util"
	"time"
)

func (ah AddressImportHandler) ParseAddressFile(fileName, datasetName string) {

	startTime := time.Now()

	ah.fileUploads.SetUploadStarted()

	rows, err := csv.ImportCSVToSlice(fileName)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "parseAddressFile").
			Str("file", fileName).
			Msg("Cannot import CSV file")
		ah.fileUploads.SetUploadError(fmt.Sprintf("cannot import CSV file %s", err))
		return
	}

	if len(rows) < 2 {
		log.Warn().
			Str("method", "parseAddressFile").
			Msg("The CSV file is empty")
		ah.fileUploads.SetUploadError("CSV file is empty import CSV file")
		return
	}

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot connect to database")
		ah.fileUploads.SetUploadError(fmt.Sprintf("cannot connect to database: %s", err))
		return
	}

	if err := database.PersistAddressDataset(rows[0], rows[1:], ah.fileUploads); err != nil {
		log.Error().
			Err(err).
			Str("datasetName", datasetName).
			Msg("Cannot persist dataset")
	}

	log.Debug().
		Str("datasetName", datasetName).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Imported and persisted dataset")

}
