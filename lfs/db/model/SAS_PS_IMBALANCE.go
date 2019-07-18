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

type SASPSIMBALANCE struct {
	FLOW       sql.NullFloat64 `gorm:"column:FLOW;primary_key"`
	SUMPRIORWT sql.NullFloat64 `gorm:"column:SUM_PRIOR_WT"`
	SUMIMBALWT sql.NullFloat64 `gorm:"column:SUM_IMBAL_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSIMBALANCE) TableName() string {
	return "SAS_PS_IMBALANCE"
}
