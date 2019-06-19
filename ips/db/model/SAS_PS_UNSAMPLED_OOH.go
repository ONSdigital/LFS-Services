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

type SASPSUNSAMPLEDOOH struct {
	UNSAMPPORTGRPPV    string          `gorm:"column:UNSAMP_PORT_GRP_PV;primary_key"`
	ARRIVEDEPART       float64         `gorm:"column:ARRIVEDEPART"`
	UNSAMPREGIONGRPPV  sql.NullFloat64 `gorm:"column:UNSAMP_REGION_GRP_PV"`
	CASES              sql.NullFloat64 `gorm:"column:CASES"`
	SUMPRIORWT         sql.NullFloat64 `gorm:"column:SUM_PRIOR_WT"`
	SUMUNSAMPTRAFFICWT sql.NullFloat64 `gorm:"column:SUM_UNSAMP_TRAFFIC_WT"`
	UNSAMPTRAFFICWT    sql.NullFloat64 `gorm:"column:UNSAMP_TRAFFIC_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSUNSAMPLEDOOH) TableName() string {
	return "SAS_PS_UNSAMPLED_OOH"
}
