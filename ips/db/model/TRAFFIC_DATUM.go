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

type TRAFFICDATUM struct {
	RUNID        string          `gorm:"column:RUN_ID;primary_key"`
	YEAR         float64         `gorm:"column:YEAR"`
	MONTH        float64         `gorm:"column:MONTH"`
	DATASOURCEID float64         `gorm:"column:DATA_SOURCE_ID"`
	PORTROUTE    sql.NullFloat64 `gorm:"column:PORTROUTE"`
	ARRIVEDEPART sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	TRAFFICTOTAL float64         `gorm:"column:TRAFFICTOTAL"`
	PERIODSTART  sql.NullString  `gorm:"column:PERIODSTART"`
	PERIODEND    sql.NullString  `gorm:"column:PERIODEND"`
	AMPMNIGHT    sql.NullFloat64 `gorm:"column:AM_PM_NIGHT"`
	HAUL         sql.NullString  `gorm:"column:HAUL"`
	VEHICLE      sql.NullFloat64 `gorm:"column:VEHICLE"`
}

// TableName sets the insert table name for this struct type
func (t *TRAFFICDATUM) TableName() string {
	return "TRAFFIC_DATA"
}
