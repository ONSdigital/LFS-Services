package db

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/db/postgres"
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

	cachedConnection = &postgres.Postgres{nil}

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

	// Import
	PersistSurvey(vo types.SurveyVO) error
	PersistAddressDataset(headers []string, rows [][]string, status *types.WSMessage) error
	GetUserID(user string) (types.UserCredentials, error)

	// Batch
	MonthlyBatchExists(month, year int) bool
	SuccessfulMonthlyBatchesExist(year int) bool
	SuccessfulQuarterlyBatchesExist(year int) bool
	AnnualBatchExists(year int) bool

	CreateMonthlyBatch(batch types.MonthlyBatch) error
	CreateAnnualBatch(batch types.AnnualBatch) error

	FindGBBatchInfo(week, year int) (types.GBBatchItem, error)
	FindNIBatchInfo(month, year int) (types.NIBatchItem, error)

	// Audits
	GetAllAudits() ([]types.Audit, error)
	GetAuditsByYear(year types.Year) ([]types.Audit, error)
	GetAuditsByYearMonth(month types.Month, year types.Year) ([]types.Audit, error)
	GetAuditsByYearWeek(week types.Week, year types.Year) ([]types.Audit, error)
}
