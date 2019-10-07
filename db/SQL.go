package db

import (
	logger "github.com/sirupsen/logrus"
	"services/dataset"
	"sync"
)

type Connection struct {
	Connection Persistence
	Log        *logger.Logger
	GlobalLock sync.Mutex
	Server     string
	UserId     string
	Password   string
	DBName     string
	Verbose    bool
}

func GetDefaultPersistenceImpl(log *logger.Logger) Persistence {
	// Maybe get this from config
	return &MySQL{nil, log}
}

type Persistence interface {
	Connect() error
	Close()
	PersistDataset(d dataset.Dataset) error
	UnpersistDataset(tableName string) (dataset.Dataset, error)
}
