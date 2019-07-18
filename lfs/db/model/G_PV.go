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

type GPV struct {
	PVID int    `gorm:"column:PV_ID;primary_key"`
	Name string `gorm:"column:Name"`
}

// TableName sets the insert table name for this struct type
func (g *GPV) TableName() string {
	return "G_PVs"
}
