package db

import (
	logger "github.com/sirupsen/logrus"
	"pds-go/lfs/db/mysql"
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

func (c Connection) DropTable(name string) error {
	return c.Connection.DropTable(name)
}

func (c Connection) Close() {
	c.Connection.Close()
}

func GetPersistenceImpl() Persistence {
	// Maybe get this from config
	return mysql.MySQL{}
}

type Persistence interface {
	Connect(*logger.Logger) error
	DropTable(name string) error
	Close()
	Insert(name string) error
}
