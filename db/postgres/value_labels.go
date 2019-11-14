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
		panic("value labels table configuration not set")
	}
}

func (s Postgres) PersistValues(d types.ValueLabelsRow) error {

	col := s.DB.Collection(valueLabelsTable)
	_, err := col.Insert(d)
	if err != nil {
		return err
	}

	return nil
}

func (s Postgres) getAllGBValueLabels() ([]types.ValueLabelsRow, error) {

	var valueLabels []types.ValueLabelsRow
	res := s.DB.Collection(valueLabelsTable).Find(db.Cond{"source": string(types.GBSource)})
	err := res.All(&valueLabels)
	if err != nil {
		return nil, res.Err()
	}

	return valueLabels, nil
}

func (s Postgres) getAllNIValueLabels() ([]types.ValueLabelsRow, error) {

	var valueLabels []types.ValueLabelsRow
	res := s.DB.Collection(valueLabelsTable).Find(db.Cond{"source": string(types.NISource)})
	err := res.All(&valueLabels)
	if err != nil {
		return nil, res.Err()
	}

	return valueLabels, nil
}

func (s Postgres) GetAllValueLabels() ([]types.ValueLabelsRow, error) {

	var valueLabels []types.ValueLabelsRow
	res := s.DB.Collection(valueLabelsTable).Find()
	err := res.All(&valueLabels)
	if err != nil {
		return nil, res.Err()
	}

	return valueLabels, nil
}

func (s Postgres) GetLabelsForValue(value string) ([]types.ValueLabelsRow, error) {

	var values []types.ValueLabelsRow
	res := s.DB.Collection(valueLabelsTable).Find(db.Cond{"name": value})

	err := res.All(&values)
	if err != nil {
		return nil, res.Err()
	}

	return values, nil
}

/* persist any new or changed value labels
 */
func (s Postgres) PersistValueLabels(header []types.ValueLabelsRow) error {

	// get existing items
	var all []types.ValueLabelsRow
	var err error
	all, err = s.GetAllValueLabels()

	if err != nil {
		return err
	}

	var newItems = make(map[string]types.ValueLabelsRow)
	for _, v := range all {
		newItems[v.Name] = v
	}

	changes := make([]types.ValueLabelsRow, 0)

	for _, v := range header {
		item, ok := newItems[v.Name]
		if !ok || item.Label != v.Label {

			r := types.ValueLabelsRow{
				Name:         v.Name,
				Label:        v.Label,
				Source:       v.Source,
				VariableType: v.VariableType,
			}
			changes = append(changes, r)
		}
	}

	if len(changes) > 0 {
		return s.persistValLabChanges(changes)
	}

	log.Info().Msg("No new or changed value labels")

	return nil
}

func (s Postgres) persistValLabChanges(values []types.ValueLabelsRow) error {

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
			Str("value", j.Name).
			Str("label", j.Label).
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
		Msg("Persisted new or changed value labels")

	return nil
}
