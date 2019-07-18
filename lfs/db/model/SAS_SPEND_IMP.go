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

type SASSPENDIMP struct {
	SERIAL   float64         `gorm:"column:SERIAL;primary_key"`
	NEWSPEND sql.NullFloat64 `gorm:"column:NEWSPEND"`
	SPENDK   sql.NullFloat64 `gorm:"column:SPENDK"`
}

// TableName sets the insert table name for this struct type
func (s *SASSPENDIMP) TableName() string {
	return "SAS_SPEND_IMP"
}
