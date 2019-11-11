package types

type VariableDefinitions struct {
	Id             int     `db:"id,omitempty"`
	Variable       string  `db:"variable" json:"variable"`
	Description    string  `db:"description" json:"description"`
	VariableType   SavType `db:"type" json:"type"`
	VariableLength int     `db:"length "json:"length"`
	Precision      int     `db:"precision" json:"precision"`
	Alias          string  `db:"alias" json:"alias"`
	Editable       bool    `db:"editable" json:"editable"`
	Imputation     bool    `db:"imputation" json:"imputation"`
	DV             bool    `db:"dv" json:"dv"`
}
