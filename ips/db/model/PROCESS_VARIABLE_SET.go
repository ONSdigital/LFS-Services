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

type PROCESSVARIABLESET struct {
	RUNID  string         `gorm:"column:RUN_ID;primary_key"`
	NAME   string         `gorm:"column:NAME"`
	USER   sql.NullString `gorm:"column:USER"`
	PERIOD string         `gorm:"column:PERIOD"`
}

// TableName sets the insert table name for this struct type
func (p *PROCESSVARIABLESET) TableName() string {
	return "PROCESS_VARIABLE_SET"
}
