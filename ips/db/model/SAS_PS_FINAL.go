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

type SASPSFINAL struct {
	SERIAL          float64         `gorm:"column:SERIAL;primary_key"`
	SHIFTWT         sql.NullFloat64 `gorm:"column:SHIFT_WT"`
	NONRESPONSEWT   sql.NullFloat64 `gorm:"column:NON_RESPONSE_WT"`
	MINSWT          sql.NullFloat64 `gorm:"column:MINS_WT"`
	TRAFFICWT       sql.NullFloat64 `gorm:"column:TRAFFIC_WT"`
	UNSAMPTRAFFICWT sql.NullFloat64 `gorm:"column:UNSAMP_TRAFFIC_WT"`
	IMBALWT         sql.NullFloat64 `gorm:"column:IMBAL_WT"`
	FINALWT         sql.NullFloat64 `gorm:"column:FINAL_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSFINAL) TableName() string {
	return "SAS_PS_FINAL"
}
