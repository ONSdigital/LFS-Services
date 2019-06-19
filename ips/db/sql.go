package db

import (
	"pds-go/ips/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

var DB *gorm.DB

func init() {
	log.Info("initialising DB")

	server := config.DBServer
	user := config.DBUser
	pass := config.DBPassword
	dbName := config.DBName
	verbose := config.Config.Database.Verbose

	connectionString := user + ":" + pass + "@tcp(" + server + ")/" + dbName + "?charset=utf8&parseTime=True"

	log.Debug("Connecting to database: ", dbName)

	var err error
	DB, err = gorm.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	log.Debug(fmt.Sprintf("Connected to database: %s", dbName ))

	if verbose {
		DB.LogMode(true)
	}

	conn := DB.DB()

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	poolSize := config.Config.Database.ConnectionPool.MaxPoolSize
	maxIdle := config.Config.Database.ConnectionPool.MaxIdleConnections
	maxLifetime := config.Config.Database.ConnectionPool.MaxLifetimeSeconds

	if maxLifetime > 0 {
		maxLifetime = maxLifetime * time.Second
		conn.SetConnMaxLifetime(maxLifetime)
		log.Debug("MaxLifetime: ", maxLifetime)
	}

	log.Debug("MaxPoolSize: ", poolSize)
	log.Debug("MaxIdleConnections: ", maxIdle)

	conn.SetMaxOpenConns(poolSize)
	conn.SetMaxIdleConns(maxIdle)

}
