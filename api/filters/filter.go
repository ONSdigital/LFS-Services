package filters

type Filter interface {
	Validate() error
	RenameColumns()
	DropColumns()
}
