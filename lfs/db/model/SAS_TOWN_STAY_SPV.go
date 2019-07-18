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

type SASTOWNSTAYSPV struct {
	SERIAL              float64         `gorm:"column:SERIAL;primary_key"`
	PURPOSEPV           sql.NullFloat64 `gorm:"column:PURPOSE_PV"`
	STAYIMPCTRYLEVEL1PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL1_PV"`
	STAYIMPCTRYLEVEL2PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL2_PV"`
	STAYIMPCTRYLEVEL3PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL3_PV"`
	STAYIMPCTRYLEVEL4PV sql.NullFloat64 `gorm:"column:STAYIMPCTRYLEVEL4_PV"`
	TOWNIMPELIGIBLEPV   sql.NullFloat64 `gorm:"column:TOWN_IMP_ELIGIBLE_PV"`
}

// TableName sets the insert table name for this struct type
func (s *SASTOWNSTAYSPV) TableName() string {
	return "SAS_TOWN_STAY_SPV"
}
