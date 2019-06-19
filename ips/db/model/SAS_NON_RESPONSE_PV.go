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

type SASNONRESPONSEPV struct {
	RECID        float64         `gorm:"column:REC_ID;primary_key"`
	WEEKDAYENDPV sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
	NRPORTGRPPV  sql.NullFloat64 `gorm:"column:NR_PORT_GRP_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASNONRESPONSEPV) TableName() string {
	return "SAS_NON_RESPONSE_PV"
}
