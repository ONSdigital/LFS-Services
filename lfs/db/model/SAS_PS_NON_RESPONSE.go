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

type SASPSNONRESPONSE struct {
	NRPORTGRPPV   float64         `gorm:"column:NR_PORT_GRP_PV;primary_key"`
	ARRIVEDEPART  float64         `gorm:"column:ARRIVEDEPART"`
	WEEKDAYENDPV  sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
	MEANRESPSSHWT sql.NullFloat64 `gorm:"column:MEAN_RESPS_SH_WT"`
	COUNTRESPS    sql.NullFloat64 `gorm:"column:COUNT_RESPS"`
	PRIORSUM      sql.NullFloat64 `gorm:"column:PRIOR_SUM"`
	GROSSRESP     sql.NullFloat64 `gorm:"column:GROSS_RESP"`
	GNR           sql.NullFloat64 `gorm:"column:GNR"`
	MEANNRWT      sql.NullFloat64 `gorm:"column:MEAN_NR_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSNONRESPONSE) TableName() string {
	return "SAS_PS_NON_RESPONSE"
}
