package api

import (
	"github.com/rs/zerolog/log"
	"services/db"
	"services/types"
)

func (vd VariableDefinitionsHandler) getAllVD() ([]types.VariableDefinitions, error) {
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

func (vd VariableDefinitionsHandler) getVDByVariable(variable string) ([]types.VariableDefinitions, error) {
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
