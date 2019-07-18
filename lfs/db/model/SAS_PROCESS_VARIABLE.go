package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type SASPROCESSVARIABLE struct {
	PROCVARNAME  string  `gorm:"column:PROCVAR_NAME;primary_key"`
	PROCVARRULE  string  `gorm:"column:PROCVAR_RULE"`
	PROCVARORDER float64 `gorm:"column:PROCVAR_ORDER"`
}

// TableName sets the insert table name for this struct type
func (s *SASPROCESSVARIABLE) TableName() string {
	return "SAS_PROCESS_VARIABLE"
}
