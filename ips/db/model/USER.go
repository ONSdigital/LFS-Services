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

type USER struct {
	ID        int            `gorm:"column:ID;primary_key"`
	USERNAME  sql.NullString `gorm:"column:USER_NAME"`
	PASSWORD  sql.NullString `gorm:"column:PASSWORD"`
	FIRSTNAME sql.NullString `gorm:"column:FIRST_NAME"`
	SURNAME   sql.NullString `gorm:"column:SURNAME"`
	ROLE      sql.NullString `gorm:"column:ROLE"`
}

// TableName sets the insert table name for this struct type
func (u *USER) TableName() string {
	return "USER"
}
