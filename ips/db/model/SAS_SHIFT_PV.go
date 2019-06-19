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

type SASSHIFTPV struct {
	RECID          float64         `gorm:"column:REC_ID;primary_key"`
	SHIFTPORTGRPPV sql.NullString  `gorm:"column:SHIFT_PORT_GRP_PV"`
	AMPMNIGHTPV    sql.NullFloat64 `gorm:"column:AM_PM_NIGHT_PV"`
	WEEKDAYENDPV   sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASSHIFTPV) TableName() string {
	return "SAS_SHIFT_PV"
}
