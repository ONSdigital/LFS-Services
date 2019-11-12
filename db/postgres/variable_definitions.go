package postgres

import (
	"fmt"
	"github.com/rs/zerolog/log"
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

func (s Postgres) PersistDefinitions(d types.VariableDefinitions) error {

	col := s.DB.Collection(definitionsTable)
	_, err := col.Insert(d)
	if err != nil {
		return err
	}

	return nil
}

func (s Postgres) GetAllDefinitions() ([]types.VariableDefinitions, error) {

	var definitions []types.VariableDefinitions
	res := s.DB.Collection(definitionsTable).Find()
	err := res.All(&definitions)
	if err != nil {
		return nil, res.Err()
	}

	return definitions, nil
}

func (s Postgres) GetDefinitionsForVariable(variable string) ([]types.VariableDefinitions, error) {

	var definitions []types.VariableDefinitions
	res := s.DB.Collection(definitionsTable).Find(db.Cond{"variable": variable})

	err := res.All(&definitions)
	if err != nil {
		return nil, res.Err()
	}

	return definitions, nil
}

/* persist any new variable definitions.
New is defined as any changes to the description
*/
func (s Postgres) PersistVariableDefinitions(header []types.Header) error {

	// get existing items
	all, err := s.GetAllDefinitions()
	if err != nil {
		return err
	}

	var newItems = make(map[string]types.VariableDefinitions)
	for _, v := range all {
		newItems[v.Variable] = v
	}

	changes := make([]types.VariableDefinitions, 0)

	for _, v := range header {
		item, ok := newItems[v.VariableName]
		if !ok || item.Description != v.VariableDescription {

			r := types.VariableDefinitions{
				Variable:       v.VariableName,
				Description:    v.VariableDescription,
				VariableType:   v.VariableType,
				VariableLength: v.VariableLength,
				Precision:      v.VariablePrecision,
				Alias:          "",
				Editable:       false,
				Imputation:     false,
				DV:             false,
			}
			changes = append(changes, r)
		}
	}

	if len(changes) > 0 {
		return s.PersistDVChanges(changes)
	} else {
		log.Info().Msg("No new or changed variable definitions")
	}

	return nil
}

func (s Postgres) PersistDVChanges(definitions []types.VariableDefinitions) error {

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	col := s.DB.Collection(definitionsTable)

	for _, j := range definitions {
		_, err = col.Insert(j)
		if err != nil {
			_ = tx.Rollback()
			log.Error().
				Err(err).
				Msg("insert into variable_definitions failed")
			return fmt.Errorf("insert into variable_definitions failed, error: %s", err)
		}
		log.Debug().
			Str("variable", j.Variable).
			Msg("Inserted variable definition")
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	log.Info().
		Int("numberItems", len(definitions)).
		Msg("Persisted new or changed variable definitions")

	return nil
}
