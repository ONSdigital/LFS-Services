package config

type configuration struct {
	Debug         bool
	TestDirectory string
	Database      DatabaseConfiguration
	Service       ServiceConfiguration
	Rename        Rename
	DropColumns   DropColumns
}
