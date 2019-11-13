package types

type VariableDefinitions struct {
	Id             int     `db:"id,omitempty"`
	Variable       string  `db:"variable" json:"variable"`
	Source         string  `db:"source" json:"source"`
	Description    string  `db:"description" json:"description"`
	VariableType   SavType `db:"type" json:"type"`
	VariableLength int     `db:"length "json:"length"`
	Precision      int     `db:"precision" json:"precision"`
	Alias          string  `db:"alias" json:"alias"`
	Editable       bool    `db:"editable" json:"editable"`
	Imputation     bool    `db:"imputation" json:"imputation"`
	DV             bool    `db:"dv" json:"dv"`
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
