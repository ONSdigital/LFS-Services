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
	headers    []string
	data       [][]string
	validation validate.Validation
	filter     Filter
	StructType interface{}
	audit      *types.Audit
	surveyType types.FileOrigin
}

func NewNIPipeLine(headers []string, rows [][]string, audit *types.Audit) Pipeline {
	return Pipeline{
		headers:    headers,
		data:       rows,
		validation: validate.NewNISurveyValidation(headers, rows),
		filter:     NewNISurveyFilter(audit),
		StructType: types.NISurveyInput{},
		audit:      audit,
		surveyType: types.NI,
	}
}

func NewGBPipeLine(headers []string, rows [][]string, audit *types.Audit) Pipeline {
	return Pipeline{
		headers:    headers,
		data:       rows,
		validation: validate.NewGBSurveyValidation(headers, rows),
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

	// Skip rows
	p.data, err = p.filter.SkipRowsFilter(p.headers, p.data)
	if err != nil {
		return nil, nil, err
	}

	// add variables
	newColumns, err := p.filter.AddVariables(p.headers, p.data)
	if err != nil {
		log.Error().
			Err(err)
		return nil, nil, err
	}

	// rename variables
	for k, v := range p.headers {
		to, ok := p.filter.RenameColumns(v)
		if ok {
			p.headers[k] = to
		}
	}

	t1 := reflect.TypeOf(p.StructType)
	columns := make([]types.Column, len(p.headers))

	colNo := 1
	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		col := types.Column{}
		// skip columns that are marked as being dropped
		if p.filter.DropColumn(strings.ToUpper(a.Name)) {
			col.Skip = true
			columns[i] = col
			continue
		}

		col.Skip = false
		col.Kind = a.Type.Kind()
		col.Name = a.Name
		col.ColNo = colNo
		colNo++
		columns[i] = col
	}

	columns = append(columns, newColumns...)

	return columns, p.data, nil
}
