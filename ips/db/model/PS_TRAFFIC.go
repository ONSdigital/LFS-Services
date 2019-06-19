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

type PSTRAFFIC struct {
	RUNID           string          `gorm:"column:RUN_ID;primary_key"`
	SAMPPORTGRPPV   string          `gorm:"column:SAMP_PORT_GRP_PV"`
	ARRIVEDEPART    float64         `gorm:"column:ARRIVEDEPART"`
	FOOTORVEHICLEPV sql.NullFloat64 `gorm:"column:FOOT_OR_VEHICLE_PV"`
	CASES           sql.NullFloat64 `gorm:"column:CASES"`
	TRAFFICTOTAL    sql.NullFloat64 `gorm:"column:TRAFFICTOTAL"`
	SUMTRAFFICWT    sql.NullFloat64 `gorm:"column:SUM_TRAFFIC_WT"`
	TRAFFICWT       sql.NullFloat64 `gorm:"column:TRAFFIC_WT"`
}

// TableName sets the insert table name for this struct type
func (p *PSTRAFFIC) TableName() string {
	return "PS_TRAFFIC"
}
