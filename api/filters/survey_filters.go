package filters

import (
	"fmt"
	conf "services/config"
	"services/dataset"
)

func NewSurveyFilter(dataset *dataset.Dataset) *DropSurveyColumns {
	return &DropSurveyColumns{dataset: dataset}
}

type DropSurveyColumns struct {
	dataset *dataset.Dataset
}

func (ren DropSurveyColumns) DropColumns() {
	drop := conf.Config.DropColumns.Survey

	for _, v := range drop.ColumnNames {
		fmt.Printf("Drop Column: %s\n", v)
	}
}

func (ren DropSurveyColumns) RenameColumns() {
	cols := conf.Config.Rename.Survey

	for _, v := range cols {
		fmt.Printf("From: %s, to: %s\n", v.From, v.To)
	}
}
