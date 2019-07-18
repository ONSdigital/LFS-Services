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

type RTRAFFIC struct {
	Rownames         sql.NullString  `gorm:"column:rownames;primary_key"`
	SERIAL           sql.NullFloat64 `gorm:"column:SERIAL"`
	ARRIVEDEPART     sql.NullInt64   `gorm:"column:ARRIVEDEPART"`
	PORTROUTE        sql.NullInt64   `gorm:"column:PORTROUTE"`
	SAMPPORTGRPPV    sql.NullString  `gorm:"column:SAMP_PORT_GRP_PV"`
	SHIFTWT          sql.NullFloat64 `gorm:"column:SHIFT_WT"`
	NONRESPONSEWT    sql.NullFloat64 `gorm:"column:NON_RESPONSE_WT"`
	MINSWT           sql.NullFloat64 `gorm:"column:MINS_WT"`
	TRAFFICWT        sql.NullFloat64 `gorm:"column:TRAFFIC_WT"`
	TRAFDESIGNWEIGHT sql.NullFloat64 `gorm:"column:TRAF_DESIGN_WEIGHT"`
	T1               sql.NullInt64   `gorm:"column:T1"`
	T                sql.NullString  `gorm:"column:T_"`
	TWWEIGHT         sql.NullFloat64 `gorm:"column:TW_WEIGHT"`
}

// TableName sets the insert table name for this struct type
func (r *RTRAFFIC) TableName() string {
	return "R_TRAFFIC"
}
