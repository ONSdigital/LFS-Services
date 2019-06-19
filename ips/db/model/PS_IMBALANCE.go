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

type PSIMBALANCE struct {
	RUNID      string          `gorm:"column:RUN_ID;primary_key"`
	FLOW       sql.NullFloat64 `gorm:"column:FLOW"`
	SUMPRIORWT sql.NullFloat64 `gorm:"column:SUM_PRIOR_WT"`
	SUMIMBALWT sql.NullFloat64 `gorm:"column:SUM_IMBAL_WT"`
}

// TableName sets the insert table name for this struct type
func (p *PSIMBALANCE) TableName() string {
	return "PS_IMBALANCE"
}
