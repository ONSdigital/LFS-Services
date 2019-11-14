package postgres

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
	"time"
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

func (s Postgres) PersistSavValueLabels(items map[string][]types.Labels, source types.FileSource) error {

	var all = make([]types.ValueLabelsRow, 0, len(items))

	for _, v := range items {
		for _, j := range v {

			item := types.ValueLabelsRow{
				Id:           0,
				Name:         j.Name,
				Label:        j.Label,
				Source:       string(source),
				VariableType: j.VariableType,
				LastUpdated:  time.Now(),
			}

			switch j.VariableType {
			case types.TypeString:
				item.Value = fmt.Sprintf("%s", j.Value)
			case types.TypeInt8, types.TypeInt16, types.TypeInt32:
				item.Value = fmt.Sprintf("%d", j.Value)
			case types.TypeFloat, types.TypeDouble:
				item.Value = fmt.Sprintf("%f", j.Value)
			}

			all = append(all, item)
		}
	}
	return s.PersistValueLabels(all)
}

/* persist any new or changed value labels
 */
func (s Postgres) PersistValueLabels(data []types.ValueLabelsRow) error {

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

	for _, v := range data {
		item, ok := newItems[v.Name]
		if !ok || (item.Label != v.Label && item.Name != v.Name && item.Value != v.Value) {

			r := types.ValueLabelsRow{
				Name:         v.Name,
				Label:        v.Label,
				Source:       v.Source,
				VariableType: v.VariableType,
				Value:        v.Value,
				LastUpdated:  v.LastUpdated,
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
			Str("name", j.Name).
			Str("label", j.Label).
			Str("value", j.Value).
			Msg("Inserted value label")
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
