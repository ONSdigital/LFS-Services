package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/importdata"
	"services/types"
	"strings"
)

func (vl ValueLabelsHandler) getAllVL() ([]types.ValueLabels, error) {
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

func (vl ValueLabelsHandler) getValLabByValue(value string) ([]types.ValueLabels, error) {
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// TODO: Persistence
	res, err := dbase.GetLabelsForValue(value)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (vl ValueLabelsHandler) parseValLabUpload(tmpfile, fileName string) error {
	var csvFile []types.ValueLabels

	if err := importdata.ImportCSVFile(tmpfile, &csvFile); err != nil {
		return err
	}

	if len(csvFile) < 1 {
		log.Warn().
			Str("method", "parseValLabUpload").
			Msg("The CSV file is empty")
		return fmt.Errorf("CSV file: %s, is empty", fileName)
	}

	// TODO: ....uuummmm.....
	v := make([]types.ValueLabels, len(csvFile))
	for i, j := range csvFile {
		varLength := intConversion(j.VariableLength)
		precision := intConversion(j.Precision)

		v[i].Variable = strings.ToUpper(j.Variable)
		v[i].Description = j.Description
		v[i].VariableType = vl.mapDataType(j.VariableType)
		v[i].VariableLength = varLength
		v[i].Precision = precision
		v[i].Alias = j.Alias
		v[i].Editable = vl.mapBool(j.Editable)
		v[i].Imputation = vl.mapBool(j.Imputation)
		v[i].DV = vl.mapBool(j.DV)
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

	if err := dbase.PersistValLabChanges(v); err != nil {
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
