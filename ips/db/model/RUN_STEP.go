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

type RUNSTEP struct {
	RUNID      string  `gorm:"column:RUN_ID;primary_key"`
	STEPNUMBER float64 `gorm:"column:STEP_NUMBER"`
	STEPNAME   string  `gorm:"column:STEP_NAME"`
	STEPSTATUS float64 `gorm:"column:STEP_STATUS"`
}

// TableName sets the insert table name for this struct type
func (r *RUNSTEP) TableName() string {
	return "RUN_STEPS"
}
