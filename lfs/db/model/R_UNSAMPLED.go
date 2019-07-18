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

type RUNSAMPLED struct {
	Rownames            sql.NullString  `gorm:"column:rownames;primary_key"`
	SERIAL              sql.NullFloat64 `gorm:"column:SERIAL"`
	ARRIVEDEPART        sql.NullInt64   `gorm:"column:ARRIVEDEPART"`
	PORTROUTE           sql.NullInt64   `gorm:"column:PORTROUTE"`
	UNSAMPPORTGRPPV     sql.NullString  `gorm:"column:UNSAMP_PORT_GRP_PV"`
	UNSAMPREGIONGRPPV   sql.NullFloat64 `gorm:"column:UNSAMP_REGION_GRP_PV"`
	SHIFTWT             sql.NullFloat64 `gorm:"column:SHIFT_WT"`
	NONRESPONSEWT       sql.NullFloat64 `gorm:"column:NON_RESPONSE_WT"`
	MINSWT              sql.NullFloat64 `gorm:"column:MINS_WT"`
	UNSAMPTRAFFICWT     sql.NullFloat64 `gorm:"column:UNSAMP_TRAFFIC_WT"`
	OOHDESIGNWEIGHT     sql.NullFloat64 `gorm:"column:OOH_DESIGN_WEIGHT"`
	T1                  sql.NullInt64   `gorm:"column:T1"`
	T                   sql.NullString  `gorm:"column:T_"`
	UNSAMPTRAFFICWEIGHT sql.NullFloat64 `gorm:"column:UNSAMP_TRAFFIC_WEIGHT"`
}

// TableName sets the insert table name for this struct type
func (r *RUNSAMPLED) TableName() string {
	return "R_UNSAMPLED"
}
