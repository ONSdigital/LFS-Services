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

type GPVVariable struct {
	PVVariableID int    `gorm:"column:PV_Variable_ID;primary_key"`
	PVID         int    `gorm:"column:PV_ID"`
	Name         string `gorm:"column:Name"`
}

// TableName sets the insert table name for this struct type
func (g *GPVVariable) TableName() string {
	return "G_PV_Variables"
}
