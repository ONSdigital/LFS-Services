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

type SASPSSHIFTDATUM struct {
	SHIFTPORTGRPPV string          `gorm:"column:SHIFT_PORT_GRP_PV;primary_key"`
	ARRIVEDEPART   float64         `gorm:"column:ARRIVEDEPART"`
	WEEKDAYENDPV   sql.NullFloat64 `gorm:"column:WEEKDAY_END_PV"`
	AMPMNIGHTPV    sql.NullFloat64 `gorm:"column:AM_PM_NIGHT_PV"`
	MIGSI          sql.NullInt64   `gorm:"column:MIGSI"`
	POSSSHIFTCROSS sql.NullFloat64 `gorm:"column:POSS_SHIFT_CROSS"`
	SAMPSHIFTCROSS sql.NullFloat64 `gorm:"column:SAMP_SHIFT_CROSS"`
	MINSHWT        sql.NullFloat64 `gorm:"column:MIN_SH_WT"`
	MEANSHWT       sql.NullFloat64 `gorm:"column:MEAN_SH_WT"`
	MAXSHWT        sql.NullFloat64 `gorm:"column:MAX_SH_WT"`
	COUNTRESPS     sql.NullFloat64 `gorm:"column:COUNT_RESPS"`
	SUMSHWT        sql.NullFloat64 `gorm:"column:SUM_SH_WT"`
}

// TableName sets the insert table name for this struct type
func (s *SASPSSHIFTDATUM) TableName() string {
	return "SAS_PS_SHIFT_DATA"
}
