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

type SASSPENDSPV struct {
	SERIAL             float64         `gorm:"column:SERIAL;primary_key"`
	SPENDIMPFLAGPV     sql.NullFloat64 `gorm:"column:SPEND_IMP_FLAG_PV"`
	SPENDIMPELIGIBLEPV sql.NullFloat64 `gorm:"column:SPEND_IMP_ELIGIBLE_PV"`
	UKOSPV             sql.NullFloat64 `gorm:"column:UK_OS_PV"`
	PUR1PV             sql.NullFloat64 `gorm:"column:PUR1_PV"`
	PUR2PV             sql.NullFloat64 `gorm:"column:PUR2_PV"`
	PUR3PV             sql.NullFloat64 `gorm:"column:PUR3_PV"`
	DUR1PV             sql.NullFloat64 `gorm:"column:DUR1_PV"`
	DUR2PV             sql.NullFloat64 `gorm:"column:DUR2_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASSPENDSPV) TableName() string {
	return "SAS_SPEND_SPV"
}
