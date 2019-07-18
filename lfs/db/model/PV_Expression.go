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

type PVExpression struct {
	ExpressionID    int `gorm:"column:Expression_ID;primary_key"`
	BlockID         int `gorm:"column:Block_ID"`
	ExpressionIndex int `gorm:"column:Expression_Index"`
}

// TableName sets the insert table name for this struct type
func (p *PVExpression) TableName() string {
	return "PV_Expression"
}
