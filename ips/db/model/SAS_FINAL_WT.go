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

type SASFINALWT struct {
	SERIAL  float64         `gorm:"column:SERIAL;primary_key"`
	FINALWT sql.NullFloat64 `gorm:"column:FINAL_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASFINALWT) TableName() string {
	return "SAS_FINAL_WT"
}
