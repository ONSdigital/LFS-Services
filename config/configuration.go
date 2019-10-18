package config

type configuration struct {
	LogFormat     string
	LogLevel      string
	TmpDirectory  string
	TestDirectory string
	Database      DatabaseConfiguration
	Service       ServiceConfiguration
	Rename        Rename
	DropColumns   DropColumns
}
