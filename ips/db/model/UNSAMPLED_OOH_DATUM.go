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

type UNSAMPLEDOOHDATUM struct {
	RUNID        string          `gorm:"column:RUN_ID;primary_key"`
	YEAR         float64         `gorm:"column:YEAR"`
	MONTH        float64         `gorm:"column:MONTH"`
	DATASOURCEID float64         `gorm:"column:DATA_SOURCE_ID"`
	PORTROUTE    sql.NullFloat64 `gorm:"column:PORTROUTE"`
	REGION       sql.NullFloat64 `gorm:"column:REGION"`
	ARRIVEDEPART sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	UNSAMPTOTAL  float64         `gorm:"column:UNSAMP_TOTAL"`
}

// TableName sets the insert table name for this struct type
func (u *UNSAMPLEDOOHDATUM) TableName() string {
	return "UNSAMPLED_OOH_DATA"
}
