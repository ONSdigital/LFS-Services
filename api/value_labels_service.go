package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/importdata"
	"services/types"
)

func (vl ValueLabelsHandler) getAllVL() ([]types.ValueLabelsRow, error) {
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	res, err := dbase.GetAllValueLabels()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (vl ValueLabelsHandler) getValLabByValue(value string) ([]types.ValueLabelsRow, error) {
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	res, err := dbase.GetLabelsForValue(value)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (vl ValueLabelsHandler) parseValLabUpload(tmpfile, fileName string) error {
	var csvFile []types.ValueLabelsRow

	if err := importdata.ImportCSVFile(tmpfile, &csvFile); err != nil {
		return err
	}

	if len(csvFile) < 1 {
		log.Warn().
			Str("method", "parseValLabUpload").
			Msg("The CSV file is empty")
		return fmt.Errorf("CSV file: %s, is empty", fileName)
	}

	log.Debug().
		Str("fileName", fileName).
		Int("rowsParsed", len(csvFile)).
		Msg("Read and parsed Value Labels file")

	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if err := dbase.PersistValueLabels(csvFile); err != nil {
		log.Error().
			Err(err).
			Str("fileName", fileName).
			Msg("Cannot persist value labels")
	}

	return nil
}

// TODO: Get the list of types from Divya
func (vl ValueLabelsHandler) mapDataType(in string) types.SavType {
	switch in {
	case "NUMBER":
		return types.TypeDouble

	default:
		return types.TypeString
	}
}

func (vl ValueLabelsHandler) mapBool(in string) bool {
	if in == "Y" {
		return true
	}
	return false
}
