package filters

import (
	conf "services/config"
	"services/dataset"
)

func NewSurveyFilter(dataset *dataset.Dataset) *SurveyFilter {
	return &SurveyFilter{dataset: dataset}
}

type SurveyFilter struct {
	dataset *dataset.Dataset
}

func (sf SurveyFilter) DropColumns() {
	drop := conf.Config.DropColumns.Survey
	_ = sf.dataset.DropColumns(drop.ColumnNames)
}

func (sf SurveyFilter) RenameColumns() {
	cols := conf.Config.Rename.Survey
	m := make(map[string]string, sf.dataset.ColumnCount)

	for _, v := range cols {
		m[v.From] = v.To
	}
	_ = sf.dataset.RenameColumns(m)
}
