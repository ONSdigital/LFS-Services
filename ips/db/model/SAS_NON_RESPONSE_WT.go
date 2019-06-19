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

type SASNONRESPONSEWT struct {
	SERIAL        float64         `gorm:"column:SERIAL;primary_key"`
	NONRESPONSEWT sql.NullFloat64 `gorm:"column:NON_RESPONSE_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASNONRESPONSEWT) TableName() string {
	return "SAS_NON_RESPONSE_WT"
}
