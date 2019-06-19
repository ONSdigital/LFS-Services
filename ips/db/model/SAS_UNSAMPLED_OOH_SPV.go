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

type SASUNSAMPLEDOOHSPV struct {
	SERIAL            float64         `gorm:"column:SERIAL;primary_key"`
	UNSAMPPORTGRPPV   sql.NullString  `gorm:"column:UNSAMP_PORT_GRP_PV"`
	UNSAMPREGIONGRPPV sql.NullFloat64 `gorm:"column:UNSAMP_REGION_GRP_PV"`
	HAULPV            sql.NullString  `gorm:"column:HAUL_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASUNSAMPLEDOOHSPV) TableName() string {
	return "SAS_UNSAMPLED_OOH_SPV"
}
