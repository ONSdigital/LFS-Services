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

type SURVEYTRAFFICAUX struct {
	SERIAL           sql.NullFloat64 `gorm:"column:SERIAL;primary_key"`
	ARRIVEDEPART     sql.NullInt64   `gorm:"column:ARRIVEDEPART"`
	PORTROUTE        sql.NullInt64   `gorm:"column:PORTROUTE"`
	SAMPPORTGRPPV    sql.NullString  `gorm:"column:SAMP_PORT_GRP_PV"`
	SHIFTWT          sql.NullFloat64 `gorm:"column:SHIFT_WT"`
	NONRESPONSEWT    sql.NullFloat64 `gorm:"column:NON_RESPONSE_WT"`
	MINSWT           sql.NullFloat64 `gorm:"column:MINS_WT"`
	TRAFFICWT        sql.NullString  `gorm:"column:TRAFFIC_WT"`
	TRAFDESIGNWEIGHT sql.NullFloat64 `gorm:"column:TRAF_DESIGN_WEIGHT"`
	T1               sql.NullInt64   `gorm:"column:T1"`
}

// TableName sets the insert table name for this struct type
func (s *SURVEYTRAFFICAUX) TableName() string {
	return "SURVEY_TRAFFIC_AUX"
}
