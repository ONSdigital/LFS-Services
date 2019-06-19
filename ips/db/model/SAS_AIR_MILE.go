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

type SASAIRMILE struct {
	SERIAL    float64         `gorm:"column:SERIAL;primary_key"`
	DIRECTLEG sql.NullFloat64 `gorm:"column:DIRECTLEG"`
	OVLEG     sql.NullFloat64 `gorm:"column:OVLEG"`
	UKLEG     sql.NullFloat64 `gorm:"column:UKLEG"`
}

// TableName sets the insert table name for this struct type
func (s *SASAIRMILE) TableName() string {
	return "SAS_AIR_MILES"
}
