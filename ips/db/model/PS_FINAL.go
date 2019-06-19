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

type PSFINAL struct {
	RUNID           string          `gorm:"column:RUN_ID;primary_key"`
	SERIAL          float64         `gorm:"column:SERIAL"`
	SHIFTWT         sql.NullFloat64 `gorm:"column:SHIFT_WT"`
	NONRESPONSEWT   sql.NullFloat64 `gorm:"column:NON_RESPONSE_WT"`
	MINSWT          sql.NullFloat64 `gorm:"column:MINS_WT"`
	TRAFFICWT       sql.NullFloat64 `gorm:"column:TRAFFIC_WT"`
	UNSAMPTRAFFICWT sql.NullFloat64 `gorm:"column:UNSAMP_TRAFFIC_WT"`
	IMBALWT         sql.NullFloat64 `gorm:"column:IMBAL_WT"`
	FINALWT         sql.NullFloat64 `gorm:"column:FINAL_WT"`
}

// TableName sets the insert table name for this struct type
func (p *PSFINAL) TableName() string {
	return "PS_FINAL"
}
