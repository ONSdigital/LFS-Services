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

type SASUNSAMPLEDOOHDATUM struct {
	RECID             int             `gorm:"column:REC_ID;primary_key"`
	PORTROUTE         sql.NullFloat64 `gorm:"column:PORTROUTE"`
	REGION            sql.NullFloat64 `gorm:"column:REGION"`
	ARRIVEDEPART      sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	UNSAMPTOTAL       sql.NullFloat64 `gorm:"column:UNSAMP_TOTAL"`
	UNSAMPPORTGRPPV   sql.NullString  `gorm:"column:UNSAMP_PORT_GRP_PV"`
	UNSAMPREGIONGRPPV sql.NullFloat64 `gorm:"column:UNSAMP_REGION_GRP_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASUNSAMPLEDOOHDATUM) TableName() string {
	return "SAS_UNSAMPLED_OOH_DATA"
}
