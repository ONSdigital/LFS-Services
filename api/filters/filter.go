package filters

type Filter interface {
	RenameColumns()
	DropColumns()
}
