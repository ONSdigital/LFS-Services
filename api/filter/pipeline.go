package filter

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/api/validate"
	"services/types"
	"strings"
)

type Pipeline struct {
	data       [][]string
	validation validate.Validation
	filter     Filter
	StructType interface{}
	audit      *types.Audit
	surveyType types.FileOrigin
}

func (p Pipeline) rows() [][]string {
	return p.data[1:]
}

func (p Pipeline) header() []string {
	return p.data[0]
}

func NewNIPipeLine(data [][]string, audit *types.Audit) Pipeline {
	return Pipeline{
		data:       data,
		validation: validate.NewNISurveyValidation(data),
		filter:     NewNISurveyFilter(audit),
		StructType: types.NISurveyInput{},
		audit:      audit,
		surveyType: types.NI,
	}
}

func NewGBPipeLine(data [][]string, audit *types.Audit) Pipeline {
	return Pipeline{
		data:       data,
		validation: validate.NewGBSurveyValidation(data),
		filter:     NewGBSurveyFilter(audit),
		StructType: types.GBSurveyInput{},
		audit:      audit,
		surveyType: types.GB,
	}
}

func (p Pipeline) RunPipeline() ([]types.Column, [][]string, error) {
	var period int
	if p.surveyType == types.GB {
		period = p.audit.Week
	} else {
		period = p.audit.Month
	}

	response, err := p.validation.Validate(period, p.audit.Year)

	if err != nil {
		return nil, nil, err
	}

	if response.ValidationResult == validate.ValidationFailed {
		return nil, nil, fmt.Errorf(response.ErrorMessage)
	}

	// Rename!!

	newColumns, err := p.filter.AddVariables(p.data)
	if err != nil {
		log.Error().
			Err(err)
		return nil, nil, err
	}

	t1 := reflect.TypeOf(p.StructType)
	columns := make([]types.Column, len(p.header()))

	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		col := types.Column{}
		// skip columns that are marked as being dropped
		if p.filter.DropColumn(strings.ToUpper(a.Name)) {
			col.Skip = true
			continue
		}

		col.Skip = false
		col.Kind = a.Type.Kind()
		col.Name = a.Name
		col.ColNo = i
		columns[i] = col
	}

	columns = append(columns, newColumns...)

	for k, v := range columns {
		to, ok := p.filter.RenameColumns(k)
		if ok {
			m[to] = v
		} else {
			m[k] = v
		}
	}

	d.Columns = m

	p.data, err = p.filter.SkipRowsFilter(p.data)
	if err != nil {
		return nil, nil, err
	}

	return columns, p.data, nil
}
