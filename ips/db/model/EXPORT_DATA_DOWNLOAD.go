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

type EXPORTDATADOWNLOAD struct {
	RUNID            string         `gorm:"column:RUN_ID;primary_key"`
	DOWNLOADABLEDATA sql.NullString `gorm:"column:DOWNLOADABLE_DATA"`
	FILENAME         sql.NullString `gorm:"column:FILENAME"`
	SOURCETABLE      sql.NullString `gorm:"column:SOURCE_TABLE"`
	DATECREATED      sql.NullString `gorm:"column:DATE_CREATED"`
}

// TableName sets the insert table name for this struct type
func (e *EXPORTDATADOWNLOAD) TableName() string {
	return "EXPORT_DATA_DOWNLOAD"
}
