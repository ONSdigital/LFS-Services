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

type RUN struct {
	RUNID        string          `gorm:"column:RUN_ID;primary_key"`
	RUNNAME      sql.NullString  `gorm:"column:RUN_NAME"`
	RUNDESC      sql.NullString  `gorm:"column:RUN_DESC"`
	USERID       sql.NullString  `gorm:"column:USER_ID"`
	PERIOD       sql.NullString  `gorm:"column:PERIOD"`
	RUNSTATUS    sql.NullFloat64 `gorm:"column:RUN_STATUS"`
	RUNTYPEID    sql.NullFloat64 `gorm:"column:RUN_TYPE_ID"`
	LASTMODIFIED time.Time       `gorm:"column:LAST_MODIFIED"`
	STEP         sql.NullString  `gorm:"column:STEP"`
	PERCENT      sql.NullInt64   `gorm:"column:PERCENT"`
}

// TableName sets the insert table name for this struct type
func (r *RUN) TableName() string {
	return "RUN"
}
