package types

import "time"

type ValueLabelsRow struct {
	Id           int       `db:"id,omitempty"`
	Name         string    `db:"name"  json:"name"`
	Label        string    `db:"label"  json:"label"`
	Value        string    `db:"value"  json:"value"`
	Source       string    `db:"source" json:"source"`
	VariableType SavType   `db:"type" json:"type"`
	LastUpdated  time.Time `db:"last_updated" json:"last_updated"`
}

type ValueLabelsImport struct {
	Variable string `csv:"VARIABLE"`
	Value    string `csv:"VALUE"`
	Label    string `csv:"LABEL"`
}
