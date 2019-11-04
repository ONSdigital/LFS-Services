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

	// Dashboard
	GetAnnualBatches() ([]types.Dashboard, error)
	GetQuarterlyBatches() ([]types.Dashboard, error)
	GetMonthlyBatches() ([]types.Dashboard, error)

	// Survey
	UnpersistSurveyDataset(tableName string) (dataset.Dataset, error)
	PersistSurveyDataset(d dataset.Dataset, vo types.SurveyVO) error

	// Address
	PersistAddressDataset(headers []string, rows [][]string, status *types.WSMessage) error

	// User
	GetUserID(user string) (types.UserCredentials, error)

	// New Batch
	MonthlyBatchExists(month, year int) bool
	AnnualBatchExists(year int) bool
	QuarterBatchExists(quarter, year int) bool

	ValidateMonthsForQuarterlyBatch(period, year int) bool
	ValidateMonthsForAnnualBatch(year int) bool
	ValidateQuartersForAnnualBatch(year int) bool

	CreateMonthlyBatch(batch types.MonthlyBatch) error
	CreateQuarterlyBatch(batch types.QuarterlyBatch) error
	CreateAnnualBatch(batch types.AnnualBatch) error

	FindGBBatchInfo(week, year int) (types.GBBatchItem, error)
	FindNIBatchInfo(month, year int) (types.NIBatchItem, error)

	// Batch IDs
	GetIdsByYear(year types.Year) ([]types.YearID, error)
	GetIdsByQuarter(year types.Year, quarter types.Quarter) ([]types.QuarterID, error)
	GetIdsByMonth(year types.Year, quarter types.Month) ([]types.MonthID, error)

	// Audits
	GetAllAudits() ([]types.Audit, error)
	GetAuditsByYear(year types.Year) ([]types.Audit, error)
	GetAuditsByYearMonth(month types.Month, year types.Year) ([]types.Audit, error)
	GetAuditsByYearWeek(week types.Week, year types.Year) ([]types.Audit, error)
}
