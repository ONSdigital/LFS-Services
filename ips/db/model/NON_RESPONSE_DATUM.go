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

type NONRESPONSEDATUM struct {
	RUNID        string          `gorm:"column:RUN_ID;primary_key"`
	YEAR         float64         `gorm:"column:YEAR"`
	MONTH        float64         `gorm:"column:MONTH"`
	DATASOURCEID float64         `gorm:"column:DATA_SOURCE_ID"`
	PORTROUTE    float64         `gorm:"column:PORTROUTE"`
	WEEKDAY      sql.NullFloat64 `gorm:"column:WEEKDAY"`
	ARRIVEDEPART sql.NullFloat64 `gorm:"column:ARRIVEDEPART"`
	AMPMNIGHT    sql.NullFloat64 `gorm:"column:AM_PM_NIGHT"`
	SAMPINTERVAL sql.NullFloat64 `gorm:"column:SAMPINTERVAL"`
	MIGTOTAL     sql.NullFloat64 `gorm:"column:MIGTOTAL"`
	ORDTOTAL     sql.NullFloat64 `gorm:"column:ORDTOTAL"`
}

// TableName sets the insert table name for this struct type
func (n *NONRESPONSEDATUM) TableName() string {
	return "NON_RESPONSE_DATA"
}
