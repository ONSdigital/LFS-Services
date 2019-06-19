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

type SASFARESIMP struct {
	SERIAL         float64         `gorm:"column:SERIAL;primary_key"`
	FARE           sql.NullFloat64 `gorm:"column:FARE"`
	FAREK          sql.NullFloat64 `gorm:"column:FAREK"`
	SPEND          sql.NullFloat64 `gorm:"column:SPEND"`
	OPERAPV        sql.NullFloat64 `gorm:"column:OPERA_PV"`
	SPENDIMPREASON sql.NullFloat64 `gorm:"column:SPENDIMPREASON"`
}

// TableName sets the insert table name for this struct type
func (s *SASFARESIMP) TableName() string {
	return "SAS_FARES_IMP"
}
