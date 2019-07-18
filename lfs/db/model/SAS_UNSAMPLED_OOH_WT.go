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

type SASUNSAMPLEDOOHWT struct {
	SERIAL          float64         `gorm:"column:SERIAL;primary_key"`
	UNSAMPTRAFFICWT sql.NullFloat64 `gorm:"column:UNSAMP_TRAFFIC_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASUNSAMPLEDOOHWT) TableName() string {
	return "SAS_UNSAMPLED_OOH_WT"
}
