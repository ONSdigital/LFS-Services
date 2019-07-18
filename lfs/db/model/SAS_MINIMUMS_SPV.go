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

type SASMINIMUMSSPV struct {
	SERIAL            float64         `gorm:"column:SERIAL;primary_key"`
	MINSPORTGRPPV     sql.NullFloat64 `gorm:"column:MINS_PORT_GRP_PV"`
	MINSCTRYGRPPV     sql.NullFloat64 `gorm:"column:MINS_CTRY_GRP_PV"`
	MINSNATGRPPV      sql.NullFloat64 `gorm:"column:MINS_NAT_GRP_PV"`
	MINSCTRYPORTGRPPV sql.NullString  `gorm:"column:MINS_CTRY_PORT_GRP_PV"`
	MINSQUALITYPV     sql.NullFloat64 `gorm:"column:MINS_QUALITY_PV"`
	MINSFLAGPV        sql.NullFloat64 `gorm:"column:MINS_FLAG_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASMINIMUMSSPV) TableName() string {
	return "SAS_MINIMUMS_SPV"
}
