package postgres

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
	"upper.io/db.v3"
)

var valueLabelsTable string

func init() {
	valueLabelsTable = config.Config.Database.ValueLabelsTable
	if valueLabelsTable == "" {
		panic("definitions table configuration not set")
	}
}

func (s Postgres) PersistValues(d types.ValueLabels) error {

	col := s.DB.Collection(valueLabelsTable)
	_, err := col.Insert(d)
	if err != nil {
		return err
	}

	return nil
}

func (s Postgres) GetAllValueLabels() ([]types.ValueLabels, error) {

	var valueLabels []types.ValueLabels
	res := s.DB.Collection(valueLabelsTable).Find()
	err := res.All(&valueLabels)
	if err != nil {
		return nil, res.Err()
	}

	return valueLabels, nil
}

func (s Postgres) GetLabelsForValue(value string) ([]types.ValueLabels, error) {

	var values []types.ValueLabels
	res := s.DB.Collection(valueLabelsTable).Find(db.Cond{"value": value})

	err := res.All(&values)
	if err != nil {
		return nil, res.Err()
	}

	return values, nil
}

/* persist any new value labelss.
New is defined as any changes to the description
//*/
// TODO: Is this required for Value Labels?
func (s Postgres) PersistValueLabels(header []types.ValueLabels) error {

	// get existing items
	all, err := s.GetAllValueLabels()
	if err != nil {
		return err
	}

	var newItems = make(map[string]types.ValueLabels)
	for _, v := range all {
		newItems[v.Variable] = v
	}

	changes := make([]types.ValueLabels, 0)

	// TODO: ....uummm....
	//for _, v := range header {
	//	item, ok := newItems[v.VariableName]
	//	if !ok || item.Description != v.VariableDescription {
	//
	//		r := types.ValueLabels{
	//			Variable:       v.VariableName,
	//			Description:    v.VariableDescription,
	//			VariableType:   v.VariableType,
	//			VariableLength: v.VariableLength,
	//			Precision:      v.VariablePrecision,
	//			Alias:          "",
	//			Editable:       false,
	//			Imputation:     false,
	//			DV:             false,
	//		}
	//		changes = append(changes, r)
	//	}
	//}

	if len(changes) > 0 {
		return s.PersistValLabChanges(changes)
	} else {
		log.Info().Msg("No new or changed value labels")
	}

	return nil
}

func (s Postgres) PersistValLabChanges(values []types.ValueLabels) error {

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	col := s.DB.Collection(valueLabelsTable)

	for _, j := range values {
		_, err = col.Insert(j)
		if err != nil {
			_ = tx.Rollback()
			log.Error().
				Err(err).
				Msg("insert into value_labels failed")
			return fmt.Errorf("insert into value_labels failed, error: %s", err)
		}
		log.Debug().
			// TODO: Will need to change j.Variable when types.ValueLabels are correctly declared
			Str("value", j.Variable).
			Msg("Inserted value labels")
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	log.Info().
		Int("numberItems", len(values)).
		Msg("Persisted new or changed variable definitions")

	return nil
}
