package types

import "time"

type ValueLabelsRow struct {
	Id           int       `db:"id,omitempty"`
	Name         string    `db:"name"  json:"name"`
	Label        string    `db:"label"  json:"label"`
	Value        int64     `db:"value"  json:"value"`
	Source       string    `db:"source" json:"source"`
	VariableType SavType   `db:"type" json:"type"`
	LastUpdated  time.Time `db:"last_updated" json:"last_updated"`
}

type ValueLabelsView struct {
	Variable         string    `db:"variable"  json:"variable"`
	Label            string    `db:"label_name"  json:"label_name"`
	Source           string    `db:"source" json:"source"`
	LabelValue       int       `db:"label_value"  json:"label_value"`
	LabelDescription SavType   `db:"label_description" json:"description"`
	LastUpdated      time.Time `db:"last_updated" json:"last_updated"`
}

type ValueLabelsImport struct {
	Variable string `csv:"VARIABLE"`
	Value    string `csv:"VALUE"`
	Label    string `csv:"LABEL"`
}
