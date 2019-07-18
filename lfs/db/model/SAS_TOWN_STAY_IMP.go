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

type SASTOWNSTAYIMP struct {
	SERIAL float64         `gorm:"column:SERIAL;primary_key"`
	SPEND1 sql.NullFloat64 `gorm:"column:SPEND1"`
	SPEND2 sql.NullFloat64 `gorm:"column:SPEND2"`
	SPEND3 sql.NullFloat64 `gorm:"column:SPEND3"`
	SPEND4 sql.NullFloat64 `gorm:"column:SPEND4"`
	SPEND5 sql.NullFloat64 `gorm:"column:SPEND5"`
	SPEND6 sql.NullFloat64 `gorm:"column:SPEND6"`
	SPEND7 sql.NullFloat64 `gorm:"column:SPEND7"`
	SPEND8 sql.NullFloat64 `gorm:"column:SPEND8"`
}

// TableName sets the insert table name for this struct type
func (s *SASTOWNSTAYIMP) TableName() string {
	return "SAS_TOWN_STAY_IMP"
}
