package config

type configuration struct {
	Debug         bool
	TestDirectory string
	Database      DatabaseConfiguration
}
