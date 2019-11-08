package types

/*
Base filter. To use this, use composition in concrete structs
*/

type Filter interface {
	DropColumn(name string) bool
	RenameColumns(column string) (string, bool)
	AddVariables(data [][]string) (int, []Column, error)
	SkipRow(row map[string]interface{}) bool
}
