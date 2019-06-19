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

type SASIMBALANCEWT struct {
	SERIAL  float64         `gorm:"column:SERIAL;primary_key"`
	IMBALWT sql.NullFloat64 `gorm:"column:IMBAL_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASIMBALANCEWT) TableName() string {
	return "SAS_IMBALANCE_WT"
}
