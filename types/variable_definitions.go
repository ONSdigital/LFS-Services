package types

import (
	"database/sql"
	"time"
)

type VariableDefinitions struct {
	Id             int            `db:"id,omitempty"`
	Variable       string         `db:"variable"`
	Label          sql.NullString `db:"label,omitempty" `
	Source         string         `db:"source" `
	Description    sql.NullString `db:"description"`
	VariableType   SavType        `db:"type"`
	VariableLength int            `db:"length"`
	Precision      int            `db:"precision"`
	Alias          sql.NullString `db:"alias" `
	Editable       bool           `db:"editable" `
	Imputation     bool           `db:"imputation"`
	DV             bool           `db:"dv" `
	ValidFrom      time.Time      `db:"valid_from"`
}

type VariableDefinitionsQuery struct {
	Variable       string    `json:"variable"`
	Label          string    `json:"label"`
	Source         string    `json:"source"`
	Description    string    `json:"description"`
	VariableType   SavType   `json:"type"`
	VariableLength int       `json:"length"`
	Precision      int       `json:"precision"`
	Alias          string    `json:"alias"`
	Editable       bool      `json:"editable"`
	Imputation     bool      `json:"imputation"`
	DV             bool      `json:"dv"`
	ValidFrom      time.Time `json:"validFrom"`
}

type VariableDefinitionsImport struct {
	Variable       string `csv:"VARIABLE"`
	Description    string `csv:"DESCRIPTION"`
	VariableType   string `csv:"DATA_TYPE"`
	VariableLength string `csv:"DATA_LENGTH"`
	Precision      string `csv:"DATA_PRECISION"`
	Alias          string `csv:"ALIAS"`
	SASFormatName  string `csv:"-"`
	Editable       string `csv:"EDITABLE"`
	Imputation     string `csv:"IMPUTATION"`
	DV             string `csv:"USED_FOR_DV"`
}
