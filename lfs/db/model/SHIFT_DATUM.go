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

type SHIFTDATUM struct {
	RUNID        string  `gorm:"column:RUN_ID;primary_key"`
	YEAR         float64 `gorm:"column:YEAR"`
	MONTH        float64 `gorm:"column:MONTH"`
	DATASOURCEID float64 `gorm:"column:DATA_SOURCE_ID"`
	PORTROUTE    float64 `gorm:"column:PORTROUTE"`
	WEEKDAY      float64 `gorm:"column:WEEKDAY"`
	ARRIVEDEPART float64 `gorm:"column:ARRIVEDEPART"`
	TOTAL        float64 `gorm:"column:TOTAL"`
	AMPMNIGHT    float64 `gorm:"column:AM_PM_NIGHT"`
}

// TableName sets the insert table name for this struct type
func (s *SHIFTDATUM) TableName() string {
	return "SHIFT_DATA"
}
