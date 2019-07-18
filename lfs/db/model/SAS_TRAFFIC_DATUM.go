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

type SASTRAFFICDATUM struct {
	RECID           int             `gorm:"column:REC_ID;primary_key"`
	PORTROUTE       sql.NullFloat64 `gorm:"column:PORTROUTE"`
	ARRIVEDEPART    sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	TRAFFICTOTAL    sql.NullFloat64 `gorm:"column:TRAFFICTOTAL"`
	PERIODSTART     sql.NullString  `gorm:"column:PERIODSTART"`
	PERIODEND       sql.NullString  `gorm:"column:PERIODEND"`
	AMPMNIGHT       sql.NullFloat64 `gorm:"column:AM_PM_NIGHT"`
	HAUL            sql.NullString  `gorm:"column:HAUL"`
	VEHICLE         sql.NullFloat64 `gorm:"column:VEHICLE"`
	SAMPPORTGRPPV   sql.NullString  `gorm:"column:SAMP_PORT_GRP_PV"`
	FOOTORVEHICLEPV sql.NullFloat64 `gorm:"column:FOOT_OR_VEHICLE_PV"`
	HAULPV          sql.NullString  `gorm:"column:HAUL_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASTRAFFICDATUM) TableName() string {
	return "SAS_TRAFFIC_DATA"
}
