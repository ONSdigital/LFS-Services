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

type SASPSMINIMUM struct {
	MINSPORTGRPPV     sql.NullFloat64 `gorm:"column:MINS_PORT_GRP_PV;primary_key"`
	ARRIVEDEPART      sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	MINSCTRYGRPPV     sql.NullFloat64 `gorm:"column:MINS_CTRY_GRP_PV"`
	MINSNATGRPPV      sql.NullFloat64 `gorm:"column:MINS_NAT_GRP_PV"`
	MINSCTRYPORTGRPPV sql.NullString  `gorm:"column:MINS_CTRY_PORT_GRP_PV"`
	MINSCASES         sql.NullFloat64 `gorm:"column:MINS_CASES"`
	FULLSCASES        sql.NullFloat64 `gorm:"column:FULLS_CASES"`
	PRIORGROSSMINS    sql.NullFloat64 `gorm:"column:PRIOR_GROSS_MINS"`
	PRIORGROSSFULLS   sql.NullFloat64 `gorm:"column:PRIOR_GROSS_FULLS"`
	PRIORGROSSALL     sql.NullFloat64 `gorm:"column:PRIOR_GROSS_ALL"`
	MINSWT            sql.NullFloat64 `gorm:"column:MINS_WT"`
	POSTSUM           sql.NullFloat64 `gorm:"column:POST_SUM"`
	CASESCARRIEDFWD   sql.NullFloat64 `gorm:"column:CASES_CARRIED_FWD"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSMINIMUM) TableName() string {
	return "SAS_PS_MINIMUMS"
}
