package config

import "time"

type Pool struct {
	MaxPoolSize        int
	MaxIdleConnections int
	MaxLifetimeSeconds time.Duration
}

type DatabaseConfiguration struct {
	Server         string
	User           string
	Password       string
	Database       string
	Verbose        bool
	ConnectionPool Pool
}
