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

type SASPSTRAFFIC struct {
	SAMPPORTGRPPV   string          `gorm:"column:SAMP_PORT_GRP_PV;primary_key"`
	ARRIVEDEPART    float64         `gorm:"column:ARRIVEDEPART"`
	FOOTORVEHICLEPV sql.NullFloat64 `gorm:"column:FOOT_OR_VEHICLE_PV"`
	CASES           sql.NullFloat64 `gorm:"column:CASES"`
	TRAFFICTOTAL    sql.NullFloat64 `gorm:"column:TRAFFICTOTAL"`
	SUMTRAFFICWT    sql.NullFloat64 `gorm:"column:SUM_TRAFFIC_WT"`
	TRAFFICWT       sql.NullFloat64 `gorm:"column:TRAFFIC_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSTRAFFIC) TableName() string {
	return "SAS_PS_TRAFFIC"
}
