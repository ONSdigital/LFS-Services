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

type PROCESSVARIABLEPY struct {
	RUNID             string  `gorm:"column:RUN_ID;primary_key"`
	PROCESSVARIABLEID float64 `gorm:"column:PROCESS_VARIABLE_ID"`
	PVNAME            string  `gorm:"column:PV_NAME"`
	PVDESC            string  `gorm:"column:PV_DESC"`
	PVDEF             string  `gorm:"column:PV_DEF"`
}

// TableName sets the insert table name for this struct type
func (p *PROCESSVARIABLEPY) TableName() string {
	return "PROCESS_VARIABLE_PY"
}
