package mysql

import (
	"github.com/rs/zerolog/log"
	"services/config"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type MySQL struct {
	DB sqlbuilder.Database
}

func (s *MySQL) Connect() error {

	var settings = mysql.ConnectionURL{
		Database: config.Config.Database.Database,
		Host:     config.Config.Database.Server,
		User:     config.Config.Database.User,
		Password: config.Config.Database.Password,
	}

	log.Debug().
		Str("databaseName", config.Config.Database.Database).
		Msg("Connecting to database")

	sess, err := mysql.Open(settings)

	if err != nil {
		log.Error().
			Err(err).
			Str("databaseName", config.Config.Database.Database).
			Msg("Cannot connect to database")
		return err
	}

	log.Debug().
		Str("databaseName", config.Config.Database.Database).
		Msg("Connected to database")

	if config.Config.Database.Verbose {
		sess.SetLogging(true)
	}

	s.DB = sess

	poolSize := config.Config.Database.ConnectionPool.MaxPoolSize
	maxIdle := config.Config.Database.ConnectionPool.MaxIdleConnections
	maxLifetime := config.Config.Database.ConnectionPool.MaxLifetimeSeconds

	if maxLifetime > 0 {
		maxLifetime = maxLifetime * time.Second
		sess.SetConnMaxLifetime(maxLifetime)
	}

	log.Debug().
		Int("MaxPoolSize", poolSize).
		Int("MaxIdleConnections", maxIdle).
		Dur("MaxLifetime", maxLifetime*time.Second).
		Msg("Connection Attributes")

	sess.SetMaxOpenConns(poolSize)
	sess.SetMaxIdleConns(maxIdle)

	return nil
}

func (s MySQL) Close() {
	if s.DB != nil {
		_ = s.DB.Close()
	}
}
