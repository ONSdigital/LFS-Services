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

type SASTRAFFICPV struct {
	RECID           float64         `gorm:"column:REC_ID;primary_key"`
	SAMPPORTGRPPV   sql.NullString  `gorm:"column:SAMP_PORT_GRP_PV"`
	FOOTORVEHICLEPV sql.NullFloat64 `gorm:"column:FOOT_OR_VEHICLE_PV"`
	HAULPV          sql.NullString  `gorm:"column:HAUL_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASTRAFFICPV) TableName() string {
	return "SAS_TRAFFIC_PV"
}
