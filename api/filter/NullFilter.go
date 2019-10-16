package filter

import (
	"services/dataset"
	"services/types"
)

type NullFilter struct {
	UKFilter
}

func (n NullFilter) AddVariables() (int, error) {
	return 0, nil
}

func (n NullFilter) SkipRow(row map[string]interface{}) bool {
	return false
}

func NewNullFilter(dataset *dataset.Dataset) types.Filter {
	return NullFilter{UKFilter{BaseFilter{dataset: dataset}}}
}
