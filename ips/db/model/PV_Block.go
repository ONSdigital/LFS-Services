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

type PVBlock struct {
	BlockID    int    `gorm:"column:Block_ID;primary_key"`
	RunID      string `gorm:"column:Run_ID"`
	BlockIndex int    `gorm:"column:Block_Index"`
	PVID       int    `gorm:"column:PV_ID"`
}

// TableName sets the insert table name for this struct type
func (p *PVBlock) TableName() string {
	return "PV_Block"
}
