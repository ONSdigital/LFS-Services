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

type SASTRAFFICWT struct {
	SERIAL    float64         `gorm:"column:SERIAL;primary_key"`
	TRAFFICWT sql.NullFloat64 `gorm:"column:TRAFFIC_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASTRAFFICWT) TableName() string {
	return "SAS_TRAFFIC_WT"
}
