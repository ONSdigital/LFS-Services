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

type SASIMBALANCESPV struct {
	SERIAL            float64         `gorm:"column:SERIAL;primary_key"`
	IMBALPORTGRPPV    sql.NullFloat64 `gorm:"column:IMBAL_PORT_GRP_PV"`
	IMBALPORTSUBGRPPV sql.NullFloat64 `gorm:"column:IMBAL_PORT_SUBGRP_PV"`
	IMBALPORTFACTPV   sql.NullFloat64 `gorm:"column:IMBAL_PORT_FACT_PV"`
	IMBALCTRYGRPPV    sql.NullFloat64 `gorm:"column:IMBAL_CTRY_GRP_PV"`
	IMBALCTRYFACTPV   sql.NullFloat64 `gorm:"column:IMBAL_CTRY_FACT_PV"`
	IMBALELIGIBLEPV   sql.NullFloat64 `gorm:"column:IMBAL_ELIGIBLE_PV"`
	PURPOSEPV         sql.NullFloat64 `gorm:"column:PURPOSE_PV"`
	FLOWPV            sql.NullFloat64 `gorm:"column:FLOW_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASIMBALANCESPV) TableName() string {
	return "SAS_IMBALANCE_SPV"
}
