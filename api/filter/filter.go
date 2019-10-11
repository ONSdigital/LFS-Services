package filter

import (
	conf "services/config"
	"services/dataset"
)

/**
Base filter. To use this, use composition in concrete structs
*/

type BaseFilter struct {
	dataset *dataset.Dataset
}

func (bf BaseFilter) runBaseFilters() {
	bf.RenameColumns()
	bf.DropColumns()
}

func (bf BaseFilter) DropColumns() {
	drop := conf.Config.DropColumns.Survey
	_ = bf.dataset.DropColumns(drop.ColumnNames)
}

func (bf BaseFilter) RenameColumns() {
	cols := conf.Config.Rename.Survey
	m := make(map[string]string, bf.dataset.ColumnCount)

	for _, v := range cols {
		m[v.From] = v.To
	}
	_ = bf.dataset.RenameColumns(m)
}
