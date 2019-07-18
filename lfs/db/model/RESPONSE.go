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

type RESPONSE struct {
	RUNID        string         `gorm:"column:RUN_ID;primary_key"`
	STEPNUMBER   int            `gorm:"column:STEP_NUMBER"`
	RESPONSECODE int            `gorm:"column:RESPONSE_CODE"`
	MESSAGE      sql.NullString `gorm:"column:MESSAGE"`
	TIMESTAMP    time.Time      `gorm:"column:TIME_STAMP"`
}

// TableName sets the insert table name for this struct type
func (r *RESPONSE) TableName() string {
	return "RESPONSE"
}
