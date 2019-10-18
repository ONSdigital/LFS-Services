package db

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/dataset"
	"services/db/mysql"
	"services/types"
	"sync"
)

var cachedConnection Persistence
var connectionMux = &sync.Mutex{}

func GetDefaultPersistenceImpl() (Persistence, error) {
	connectionMux.Lock()
	defer connectionMux.Unlock()

	if cachedConnection != nil {
		log.Info().
			Str("databaseName", config.Config.Database.Database).
			Msg("Returning cached database connection")
		return cachedConnection, nil
	}

	cachedConnection = &mysql.MySQL{nil}

	if err := cachedConnection.Connect(); err != nil {
		log.Info().
			Err(err).
			Str("databaseName", config.Config.Database.Database).
			Msg("Cannot connect to database")
		cachedConnection = nil
		return nil, fmt.Errorf("cannot connect to database")
	}

	return cachedConnection, nil
}

type Persistence interface {
	Connect() error
	Close()
	PersistSurveyDataset(d dataset.Dataset) error
	PersistAddressDataset(tmpfile string) error
	UnpersistSurveyDataset(tableName string) (dataset.Dataset, error)
	GetUserID(user string) (types.UserCredentials, error)
}
