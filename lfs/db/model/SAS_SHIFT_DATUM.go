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

type SASSHIFTDATUM struct {
	RECID          int             `gorm:"column:REC_ID;primary_key"`
	PORTROUTE      float64         `gorm:"column:PORTROUTE"`
	WEEKDAY        float64         `gorm:"column:WEEKDAY"`
	ARRIVEDEPART   float64         `gorm:"column:ARRIVEDEPART"`
	TOTAL          float64         `gorm:"column:TOTAL"`
	AMPMNIGHT      float64         `gorm:"column:AM_PM_NIGHT"`
	SHIFTPORTGRPPV sql.NullString  `gorm:"column:SHIFT_PORT_GRP_PV"`
	AMPMNIGHTPV    sql.NullFloat64 `gorm:"column:AM_PM_NIGHT_PV"`
	WEEKDAYENDPV   sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASSHIFTDATUM) TableName() string {
	return "SAS_SHIFT_DATA"
}
