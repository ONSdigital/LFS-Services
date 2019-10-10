package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"services/config"
	"services/dataset"
	"sync"
)

var cachedConnection Persistence
var connectionMux = &sync.Mutex{}

func GetDefaultPersistenceImpl() (Persistence, error) {
	connectionMux.Lock()
	defer connectionMux.Unlock()

	if cachedConnection != nil {
		log.WithFields(log.Fields{
			"databaseName": config.Config.Database.Database,
		}).Debug("Returning cached database connection")
		return cachedConnection, nil
	}

	cachedConnection = &MySQL{nil}

	if err := cachedConnection.Connect(); err != nil {
		log.WithFields(log.Fields{
			"databaseName": config.Config.Database.Database,
			"errorMessage": err.Error(),
		}).Error("Cannot connect to database")
		cachedConnection = nil
		return nil, fmt.Errorf("cannot connect to database")
	}

	return cachedConnection, nil
}

type Persistence interface {
	Connect() error
	Close()
	PersistDataset(d dataset.Dataset) error
	UnpersistDataset(tableName string) (dataset.Dataset, error)
}
