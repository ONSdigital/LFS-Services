package filter

import (
	"github.com/rs/zerolog/log"
	conf "services/config"
	"services/dataset"
)

var dropColumns = conf.Config.DropColumns.Survey
var renameColumns map[string]string

func init() {
	cols := conf.Config.Rename.Survey
	renameColumns = make(map[string]string, len(cols))

	for _, v := range cols {
		renameColumns[v.From] = v.To
	}
}

type BaseFilter struct {
	dataset *dataset.Dataset
}

/*
Generic drop columns functionality - based on the name of columns to drop in the configuration file
*/
func (bf BaseFilter) DropColumn(name string) bool {
	for _, j := range dropColumns.ColumnNames {
		if j == name {
			log.Debug().
				Str("columnName", name).
				Msg("Dropping column")
			bf.dataset.NumVarLoaded = bf.dataset.NumVarLoaded - 1
			return true
		}
	}
	return false
}

/*
Generic rename columns functionality - based on the name of columns to drop in the configuration file
*/
func (bf BaseFilter) RenameColumns(column string) (string, bool) {
	item, ok := renameColumns[column]
	if ok {
		log.Debug().
			Str("from", column).
			Str("to", item).
			Msg("Renaming column")
		return item, true
	}
	return "", false
}
