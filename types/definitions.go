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

type SavType string

const (
	TYPE_STRING SavType = "string"
	TYPE_INT8   SavType = "int8"
	TYPE_INT16  SavType = "int16"
	TYPE_INT32  SavType = "int32"
	TYPE_FLOAT  SavType = "float"
	TYPE_DOUBLE SavType = "double"
)
