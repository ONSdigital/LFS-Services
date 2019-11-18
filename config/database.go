package config

import (
	"time"
)

type Pool struct {
	MaxPoolSize        int
	MaxIdleConnections int
	MaxLifetimeSeconds time.Duration
}

type DatabaseConfiguration struct {
	Server              string `env:"DB_SERVER"`
	User                string `env:"DB_USER"`
	Password            string `env:"DB_PASSWORD"`
	Database            string `env:"DB_DATABASE"`
	Verbose             bool
	ConnectionPool      Pool
	SurveyTable         string
	AddressesTable      string
	SurveyAuditTable    string
	BatchInfoView       string
	GbInfoView          string
	NiInfoView          string
	MonthlyBatchTable   string
	QuarterlyBatchTable string
	AnnualBatchTable    string
	GbBatchTable        string
	NiBatchTable        string
	UserTable           string
	DefinitionsTable    string
	ValueLabelsTable    string
	ValueLabelsView     string
}
