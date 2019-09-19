package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"pds-go/lfs/config"
	"sync"
	"time"
)

var globalLock = sync.Mutex{}

type SQL struct {
	DB     *gorm.DB
	logger *log.Logger
}

func NewSQL(logger *log.Logger) (*SQL, error) {

	globalLock.Lock()
	defer globalLock.Unlock()

	log.Info("initialising DB")
	server := config.Config.Database.Server
	user := config.Config.Database.User
	pass := config.Config.Database.Password
	dbName := config.Config.Database.Database
	verbose := config.Config.Database.Verbose

	connectionString := user + ":" + pass + "@tcp(" + server + ")/" + dbName + "?charset=utf8&parseTime=True"

	log.Debug("Connecting to database: ", dbName)

	var err error
	db, err := gorm.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot open database connection %v", err))
		return &SQL{}, err
	}

	log.Debug(fmt.Sprintf("Connected to database: %s", dbName))

	if verbose {
		db.LogMode(true)
	}

	conn := db.DB()

	err = conn.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot ping database %v", err))
		return &SQL{}, err
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

	return &SQL{db, logger}, nil
}
