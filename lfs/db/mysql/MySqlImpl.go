package mysql

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"pds-go/lfs/config"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type MySQL struct {
	DB sqlbuilder.Database
}

func (s MySQL) Connect(log *logger.Logger) error {

	var connectDSN = "mysql://" +
		config.Config.Database.User +
		":" +
		config.Config.Database.Password +
		"@" +
		config.Config.Database.Server +
		"/" +
		config.Config.Database.Database +
		"?charset=utf8&parseTime=True"

	settings, err := mysql.ParseURL(connectDSN)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot open database connection %v", err))
		return err
	}

	log.Debug("Connecting to database: ", config.Config.Database.Database)
	sess, err := mysql.Open(settings)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot open database connection %v", err))
		return err
	}

	log.Debug(fmt.Sprintf("Connected to database: %s", config.Config.Database.Database))

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
		log.Debug("MaxLifetime: ", maxLifetime)
	}

	log.Debug("MaxPoolSize: ", poolSize)
	log.Debug("MaxIdleConnections: ", maxIdle)

	sess.SetMaxOpenConns(poolSize)
	sess.SetMaxIdleConns(maxIdle)

	return nil
}

func (s MySQL) DropTable(name string) error {
	if _, err := s.DB.Exec("DROP TABLE IF EXISTS ", name); err != nil {
		return err
	}
	return nil
}

func (s MySQL) Close() {
	if s.DB != nil {
		_ = s.DB.Close()
	}
}

func (s MySQL) Insert(name string) error {
	return nil
}
