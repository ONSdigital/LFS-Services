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

/*
Base filter. To use this, use composition in concrete structs
*/

type Filter interface {
	DropColumn(name string) bool
	RenameColumn(column string) (string, bool)
	AddVariables() (int, error)
}

type BaseFilter struct {
	dataset *dataset.Dataset
}

/*
Generic drop columns functionality - based on the name of columns to drop in the configuration file
*/
func (sf GBSurveyFilter) DropColumn(name string) bool {
	for _, j := range dropColumns.ColumnNames {
		if j == name {
			log.Debug().
				Str("columnName", name).
				Msg("Dropping column")
			sf.dataset.NumVarLoaded = sf.dataset.NumVarLoaded - 1
			return true
		}
	}
	return false
}

/*
Generic rename columns functionality - based on the name of columns to drop in the configuration file
*/
func (bf BaseFilter) RenameColumn(column string) (string, bool) {
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
