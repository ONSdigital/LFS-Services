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

type SASREGIONALIMP struct {
	SERIAL         float64         `gorm:"column:SERIAL;primary_key"`
	VISITWT        sql.NullFloat64 `gorm:"column:VISIT_WT"`
	STAYWT         sql.NullFloat64 `gorm:"column:STAY_WT"`
	EXPENDITUREWT  sql.NullFloat64 `gorm:"column:EXPENDITURE_WT"`
	VISITWTK       sql.NullString  `gorm:"column:VISIT_WTK"`
	STAYWTK        sql.NullString  `gorm:"column:STAY_WTK"`
	EXPENDITUREWTK sql.NullString  `gorm:"column:EXPENDITURE_WTK"`
	NIGHTS1        sql.NullFloat64 `gorm:"column:NIGHTS1"`
	NIGHTS2        sql.NullFloat64 `gorm:"column:NIGHTS2"`
	NIGHTS3        sql.NullFloat64 `gorm:"column:NIGHTS3"`
	NIGHTS4        sql.NullFloat64 `gorm:"column:NIGHTS4"`
	NIGHTS5        sql.NullFloat64 `gorm:"column:NIGHTS5"`
	NIGHTS6        sql.NullFloat64 `gorm:"column:NIGHTS6"`
	NIGHTS7        sql.NullFloat64 `gorm:"column:NIGHTS7"`
	NIGHTS8        sql.NullFloat64 `gorm:"column:NIGHTS8"`
	STAY1K         sql.NullString  `gorm:"column:STAY1K"`
	STAY2K         sql.NullString  `gorm:"column:STAY2K"`
	STAY3K         sql.NullString  `gorm:"column:STAY3K"`
	STAY4K         sql.NullString  `gorm:"column:STAY4K"`
	STAY5K         sql.NullString  `gorm:"column:STAY5K"`
	STAY6K         sql.NullString  `gorm:"column:STAY6K"`
	STAY7K         sql.NullString  `gorm:"column:STAY7K"`
	STAY8K         sql.NullString  `gorm:"column:STAY8K"`
}

// TableName sets the insert table name for this struct type
func (s *SASREGIONALIMP) TableName() string {
	return "SAS_REGIONAL_IMP"
}
