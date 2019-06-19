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

type SASMINIMUMSWT struct {
	SERIAL float64         `gorm:"column:SERIAL;primary_key"`
	MINSWT sql.NullFloat64 `gorm:"column:MINS_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASMINIMUMSWT) TableName() string {
	return "SAS_MINIMUMS_WT"
}
