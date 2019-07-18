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

type PVElement struct {
	ElementID    int    `gorm:"column:Element_ID;primary_key"`
	ExpressionID int    `gorm:"column:Expression_ID"`
	Type         string `gorm:"column:type"`
	Content      string `gorm:"column:content"`
}

// TableName sets the insert table name for this struct type
func (p *PVElement) TableName() string {
	return "PV_Element"
}
