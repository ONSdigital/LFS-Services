package postgres

import (
	"services/config"
	"services/types"
	"upper.io/db.v3"
)

var definitionsTable string

func init() {
	definitionsTable = config.Config.Database.DefinitionsTable
	if definitionsTable == "" {
		panic("definitions table configuration not set")
	}
}

func (s Postgres) PersistDefinitions(d types.Definitions) error {

	col := s.DB.Collection(definitionsTable)
	_, err := col.Insert(d)
	if err != nil {
		return err
	}

	return nil
}

func (s Postgres) GetAllDefinitions() ([]types.Definitions, error) {

	var definitions []types.Definitions
	res := s.DB.Collection(definitionsTable).Find()
	err := res.All(&definitions)
	if err != nil {
		return nil, res.Err()
	}

	return definitions, nil
}

func (s Postgres) GetDefinitionsForVariable(variable string) ([]types.Definitions, error) {

	var definitions []types.Definitions
	res := s.DB.Collection(definitionsTable).Find(db.Cond{"variable": variable})

	err := res.All(&definitions)
	if err != nil {
		return nil, res.Err()
	}

	return definitions, nil
}
