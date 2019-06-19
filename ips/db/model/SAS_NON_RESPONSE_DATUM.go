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

type SASNONRESPONSEDATUM struct {
	RECID        int             `gorm:"column:REC_ID;primary_key"`
	PORTROUTE    float64         `gorm:"column:PORTROUTE"`
	WEEKDAY      sql.NullFloat64 `gorm:"column:WEEKDAY"`
	ARRIVEDEPART sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	AMPMNIGHT    sql.NullFloat64 `gorm:"column:AM_PM_NIGHT"`
	SAMPINTERVAL sql.NullFloat64 `gorm:"column:SAMPINTERVAL"`
	MIGTOTAL     sql.NullFloat64 `gorm:"column:MIGTOTAL"`
	ORDTOTAL     sql.NullFloat64 `gorm:"column:ORDTOTAL"`
	NRPORTGRPPV  sql.NullFloat64 `gorm:"column:NR_PORT_GRP_PV"`
	WEEKDAYENDPV sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
	AMPMNIGHTPV  sql.NullFloat64 `gorm:"column:AM_PM_NIGHT_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASNONRESPONSEDATUM) TableName() string {
	return "SAS_NON_RESPONSE_DATA"
}
