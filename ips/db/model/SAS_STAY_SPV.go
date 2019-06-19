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

type SASSTAYSPV struct {
	SERIAL              float64         `gorm:"column:SERIAL;primary_key"`
	STAYIMPFLAGPV       sql.NullFloat64 `gorm:"column:STAY_IMP_FLAG_PV"`
	STAYIMPELIGIBLEPV   sql.NullFloat64 `gorm:"column:STAY_IMP_ELIGIBLE_PV"`
	STAYIMPCTRYLEVEL1PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL1_PV"`
	STAYIMPCTRYLEVEL2PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL2_PV"`
	STAYIMPCTRYLEVEL3PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL3_PV"`
	STAYIMPCTRYLEVEL4PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL4_PV"`
	STAYPURPOSEGRPPV    sql.NullFloat64 `gorm:"column:STAY_PURPOSE_GRP_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASSTAYSPV) TableName() string {
	return "SAS_STAY_SPV"
}
