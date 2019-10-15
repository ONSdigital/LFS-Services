package types

//type DropFunction func(name string) bool
//type RenameFunction func(name string) (string, bool)
//type SkipRowFunction func(row map[string]interface{}) bool

/*
Base filter. To use this, use composition in concrete structs
*/

type Filter interface {
	DropColumn(name string) bool
	RenameColumns(column string) (string, bool)
	AddVariables() (int, error)
	SkipRow(row map[string]interface{}) bool
}
