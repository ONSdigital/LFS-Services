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

type SASNONRESPONSESPV struct {
	SERIAL      float64         `gorm:"column:SERIAL;primary_key"`
	NRPORTGRPPV sql.NullFloat64 `gorm:"column:NR_PORT_GRP_PV"`
	MIGFLAGPV   sql.NullFloat64 `gorm:"column:MIG_FLAG_PV"`
	NRFLAGPV    sql.NullFloat64 `gorm:"column:NR_FLAG_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASNONRESPONSESPV) TableName() string {
	return "SAS_NON_RESPONSE_SPV"
}
