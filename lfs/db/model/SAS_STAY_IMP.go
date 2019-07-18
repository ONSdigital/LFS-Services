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

type SASSTAYIMP struct {
	SERIAL float64         `gorm:"column:SERIAL;primary_key"`
	STAY   sql.NullFloat64 `gorm:"column:STAY"`
	STAYK  sql.NullFloat64 `gorm:"column:STAYK"`
}

// TableName sets the insert table name for this struct type
func (s *SASSTAYIMP) TableName() string {
	return "SAS_STAY_IMP"
}
