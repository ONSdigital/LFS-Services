package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/db"
	"services/importdata"
	"services/types"
	"strconv"
	"time"
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

func (vl ValueLabelsHandler) parseValLabUpload(tmpfile, fileName string, source types.FileSource) error {
	var csvFile []types.ValueLabelsImport

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

	var imp = make([]types.ValueLabelsRow, len(csvFile))
	for i, j := range csvFile {
		imp[i] = types.ValueLabelsRow{
			Name:         j.Variable,
			Label:        j.Label,
			Value:        j.Value,
			Source:       string(source),
			VariableType: getSource(j.Value),
			LastUpdated:  time.Now(),
		}
	}

	if err := dbase.PersistValueLabels(imp); err != nil {
		log.Error().
			Err(err).
			Str("fileName", fileName).
			Msg("Cannot persist value labels")
	}

	return nil
}

// we only consider string or float
func getSource(val string) types.SavType {

	_, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return types.TypeDouble
	}

	return types.TypeString

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
