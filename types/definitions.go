package types

type Definitions struct {
	Variable       string  `json:"variable"`
	Description    string  `json:"description"`
	VariableType   SavType `json:"type"`
	VariableLength int     `json:"length"`
	Precision      int     `json:"precision"`
	Alias          string  `json:"alias"`
	Editable       bool    `json:"editable"`
	Imputation     bool    `json:"imputation"`
	DV             bool    `json:"dv"`
}
