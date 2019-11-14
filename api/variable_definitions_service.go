package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/importdata"
	"services/types"
	"services/util"
	"strings"
)

func (vd VariableDefinitionsHandler) getAllVD() ([]types.VariableDefinitionsQuery, error) {
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	res, err := dbase.GetAllDefinitions()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (vd VariableDefinitionsHandler) getVDByVariable(variable string) ([]types.VariableDefinitionsQuery, error) {
	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	res, err := dbase.GetDefinitionsForVariable(variable)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (vd VariableDefinitionsHandler) parseVDUpload(tmpfile, fileName string) error {
	var csvFile []types.VariableDefinitionsImport

	if err := importdata.ImportCSVFile(tmpfile, &csvFile); err != nil {
		return err
	}

	if len(csvFile) < 1 {
		log.Warn().
			Str("method", "parseVDUpload").
			Msg("The CSV file is empty")
		return fmt.Errorf("CSV file: %s, is empty", fileName)
	}

	v := make([]types.VariableDefinitions, len(csvFile))
	for i, j := range csvFile {
		varLength := intConversion(j.VariableLength)
		precision := intConversion(j.Precision)

		v[i].Variable = strings.ToUpper(j.Variable)
		v[i].Description = util.ToNullString(j.Description)
		v[i].VariableType = vd.mapDataType(j.VariableType)
		v[i].VariableLength = varLength
		v[i].Precision = precision
		v[i].Alias = util.ToNullString(j.Alias)
		v[i].Editable = vd.mapBool(j.Editable)
		v[i].Imputation = vd.mapBool(j.Imputation)
		v[i].DV = vd.mapBool(j.DV)
	}

	log.Debug().
		Str("fileName", fileName).
		Int("rowsParsed", len(csvFile)).
		Msg("Read and parsed Variable Definitions file")

	dbase, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if err := dbase.PersistDVChanges(v); err != nil {
		log.Error().
			Err(err).
			Str("fileName", fileName).
			Msg("Cannot persist variable definitions")
	}

	return nil
}

// TODO: Get the list of types from Divya
func (vd VariableDefinitionsHandler) mapDataType(in string) types.SavType {
	switch in {
	case "NUMBER":
		return types.TypeDouble

	default:
		return types.TypeString
	}
}

func (vd VariableDefinitionsHandler) mapBool(in string) bool {
	if in == "Y" {
		return true
	}
	return false
}
