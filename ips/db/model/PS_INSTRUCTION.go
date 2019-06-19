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

type PSINSTRUCTION struct {
	PNID          float64 `gorm:"column:PN_ID;primary_key"`
	PSINSTRUCTION string  `gorm:"column:PS_INSTRUCTION"`
}

// TableName sets the insert table name for this struct type
func (p *PSINSTRUCTION) TableName() string {
	return "PS_INSTRUCTION"
}
