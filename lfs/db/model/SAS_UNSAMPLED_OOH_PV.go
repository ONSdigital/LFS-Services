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

type SASUNSAMPLEDOOHPV struct {
	RECID             float64         `gorm:"column:REC_ID;primary_key"`
	UNSAMPPORTGRPPV   string          `gorm:"column:UNSAMP_PORT_GRP_PV"`
	UNSAMPREGIONGRPPV sql.NullFloat64 `gorm:"column:UNSAMP_REGION_GRP_PV"`
	HAULPV            sql.NullString  `gorm:"column:HAUL_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASUNSAMPLEDOOHPV) TableName() string {
	return "SAS_UNSAMPLED_OOH_PV"
}
