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

type SASSHIFTSPV struct {
	SERIAL          float64         `gorm:"column:SERIAL;primary_key"`
	SHIFTPORTGRPPV  sql.NullString  `gorm:"column:SHIFT_PORT_GRP_PV"`
	AMPMNIGHTPV     sql.NullFloat64 `gorm:"column:AM_PM_NIGHT_PV"`
	WEEKDAYENDPV    sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
	SHIFTFLAGPV     sql.NullFloat64 `gorm:"column:SHIFT_FLAG_PV"`
	CROSSINGSFLAGPV sql.NullFloat64 `gorm:"column:CROSSINGS_FLAG_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASSHIFTSPV) TableName() string {
	return "SAS_SHIFT_SPV"
}
