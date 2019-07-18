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

type SASRAILIMP struct {
	SERIAL float64         `gorm:"column:SERIAL;primary_key"`
	SPEND  sql.NullFloat64 `gorm:"column:SPEND"`
}

// TableName sets the insert table name for this struct type
func (s *SASRAILIMP) TableName() string {
	return "SAS_RAIL_IMP"
}
