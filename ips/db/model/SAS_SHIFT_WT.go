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

type SASSHIFTWT struct {
	SERIAL  float64         `gorm:"column:SERIAL;primary_key"`
	SHIFTWT sql.NullFloat64 `gorm:"column:SHIFT_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASSHIFTWT) TableName() string {
	return "SAS_SHIFT_WT"
}
