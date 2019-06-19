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

type SASRAILSPV struct {
	SERIAL            float64         `gorm:"column:SERIAL;primary_key"`
	RAILCNTRYGRPPV    sql.NullFloat64 `gorm:"column:RAIL_CNTRY_GRP_PV"`
	RAILEXERCISEPV    sql.NullFloat64 `gorm:"column:RAIL_EXERCISE_PV"`
	RAILIMPELIGIBLEPV sql.NullFloat64 `gorm:"column:RAIL_IMP_ELIGIBLE_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASRAILSPV) TableName() string {
	return "SAS_RAIL_SPV"
}
