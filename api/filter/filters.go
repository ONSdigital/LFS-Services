package filter

type Filter interface {
	RenameColumns()
	DropColumns()
	AddVariables() error
}
