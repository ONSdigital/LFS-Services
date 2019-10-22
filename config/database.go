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
	Server           string `env:"DB_SERVER" envDefault:"localhost"`
	User             string `env:"DB_USER"`
	Password         string `env:"DB_PASSWORD"`
	Database         string `env:"DB_DATABASE"`
	Verbose          bool
	ConnectionPool   Pool
	ColumnsTable     string
	AddressesTable   string
	UploadAuditTable string
}
