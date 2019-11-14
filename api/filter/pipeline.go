package filter

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/api/validate"
	"services/importdata/sav"
	"services/types"
	"strings"
)

type Pipeline struct {
	data       types.SavImportData
	validation validate.Validation
	filter     Filter
	StructType interface{}
	audit      *types.Audit
	surveyType types.FileOrigin
}

func NewNIPipeLine(data types.SavImportData, audit *types.Audit) Pipeline {

	return Pipeline{
		data:       data,
		validation: nil,
		filter:     NewNISurveyFilter(audit),
		StructType: types.NISurveyInput{},
		audit:      audit,
		surveyType: types.NI,
	}
}

func NewGBPipeLine(data types.SavImportData, audit *types.Audit) Pipeline {
	return Pipeline{
		data:       data,
		validation: nil,
		filter:     NewGBSurveyFilter(audit),
		StructType: types.GBSurveyInput{},
		audit:      audit,
		surveyType: types.GB,
	}
}

func (p Pipeline) RunPipeline() ([]types.Column, [][]string, error) {
	var period int
	var headers []string
	var body [][]string

	if p.surveyType == types.GB {
		period = p.audit.Week
		headers, body = sav.SPSSDatatoArray(p.data)
		p.validation = validate.NewGBSurveyValidation(headers, body, p.data)
	} else {
		headers, body = sav.SPSSDatatoArray(p.data)
		p.validation = validate.NewNISurveyValidation(headers, body, p.data)
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
	data, err := p.filter.SkipRowsFilter(headers, body, p.data)
	if err != nil {
		return nil, nil, err
	}

	// add variables
	newColumns, err := p.filter.AddVariables(headers, data, p.data)
	if err != nil {
		log.Error().
			Err(err)
		return nil, nil, err
	}

	// rename variables
	for k, v := range headers {
		to, ok := p.filter.RenameColumns(v)
		if ok {
			headers[k] = to
		}
	}

	t1 := reflect.TypeOf(p.StructType)
	columns := make([]types.Column, len(headers))

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
		col.Name = headers[i]
		col.ColNo = colNo
		col.Label = p.data.Header[i].LabelName
		colNo++
		columns[i] = col
	}

	columns = append(columns, newColumns...)

	return columns, data, nil
}
