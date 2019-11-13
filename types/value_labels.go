package types

// TODO: I don't know what these are!?
type ValueLabels struct {
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

type ValueLabelsImport struct {
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
