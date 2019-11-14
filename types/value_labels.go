package types

import "time"

type ValueLabelsRow struct {
	Id           int       `db:"id,omitempty"`
	Name         string    `db:"name"  json:"type"`
	Label        string    `db:"label"  json:"label"`
	VariableType SavType   `db:"type" json:"type"`
	LastUpdated  time.Time `db:"last_updated" json:"last_updated"`
}