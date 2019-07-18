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

type SASFARESSPV struct {
	SERIAL              float64         `gorm:"column:SERIAL;primary_key"`
	FARESIMPFLAGPV      sql.NullFloat64 `gorm:"column:FARES_IMP_FLAG_PV"`
	FARESIMPELIGIBLEPV  sql.NullFloat64 `gorm:"column:FARES_IMP_ELIGIBLE_PV"`
	DISCNTPACKAGECOSTPV sql.NullFloat64 `gorm:"column:DISCNT_PACKAGE_COST_PV"`
	DISCNTF1PV          sql.NullFloat64 `gorm:"column:DISCNT_F1_PV"`
	DISCNTF2PV          sql.NullFloat64 `gorm:"column:DISCNT_F2_PV"`
	FAGEPV              sql.NullFloat64 `gorm:"column:FAGE_PV"`
	TYPEPV              sql.NullFloat64 `gorm:"column:TYPE_PV"`
	UKPORT1PV           sql.NullFloat64 `gorm:"column:UKPORT1_PV"`
	UKPORT2PV           sql.NullFloat64 `gorm:"column:UKPORT2_PV"`
	UKPORT3PV           sql.NullFloat64 `gorm:"column:UKPORT3_PV"`
	UKPORT4PV           sql.NullFloat64 `gorm:"column:UKPORT4_PV"`
	OSPORT1PV           sql.NullFloat64 `gorm:"column:OSPORT1_PV"`
	OSPORT2PV           sql.NullFloat64 `gorm:"column:OSPORT2_PV"`
	OSPORT3PV           sql.NullFloat64 `gorm:"column:OSPORT3_PV"`
	OSPORT4PV           sql.NullFloat64 `gorm:"column:OSPORT4_PV"`
	APDPV               sql.NullFloat64 `gorm:"column:APD_PV"`
	QMFAREPV            sql.NullFloat64 `gorm:"column:QMFARE_PV"`
	DUTYFREEPV          sql.NullFloat64 `gorm:"column:DUTY_FREE_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASFARESSPV) TableName() string {
	return "SAS_FARES_SPV"
}
