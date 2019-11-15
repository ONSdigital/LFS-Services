package filter

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/api/validate"
	"services/types"
	"strings"
)

type Pipeline struct {
	data       *types.SavImportData
	validation validate.Validation
	filter     Filter
	audit      *types.Audit
	surveyType types.FileOrigin
}

func NewNIPipeLine(data *types.SavImportData, audit *types.Audit) Pipeline {

	return Pipeline{
		data:       data,
		validation: nil,
		filter:     NewNISurveyFilter(),
		audit:      audit,
		surveyType: types.NI,
	}
}

func NewGBPipeLine(data *types.SavImportData, audit *types.Audit) Pipeline {
	return Pipeline{
		data:       data,
		validation: nil,
		filter:     NewGBSurveyFilter(),
		audit:      audit,
		surveyType: types.GB,
	}
}

func (p Pipeline) RunPipeline() error {
	var period int

	if p.surveyType == types.GB {
		period = p.audit.Week
		p.validation = validate.NewGBSurveyValidation(p.data)
	} else {
		p.validation = validate.NewNISurveyValidation(p.data)
		period = p.audit.Month
	}

	response, err := p.validation.Validate(period, p.audit.Year)

	if err != nil {
		return err
	}

	if response.ValidationResult == validate.ValidationFailed {
		return fmt.Errorf(response.ErrorMessage)
	}

	// Skip rows filter
	if err := p.filter.SkipRowsFilter(p.data); err != nil {
		return err
	}

	// add variables
	err = p.filter.AddVariables(p.data)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	// rename variables
	for k, v := range p.data.Header {
		to, ok := p.filter.RenameColumns(v.VariableName)
		if ok {
			p.data.Header[k].VariableName = to
		}
	}

	// drop unwanted columns
	headers := make([]types.Header, 0, p.data.HeaderCount)
	rowsToDrop := make(map[int]bool, p.data.HeaderCount)

	// mark columns of rows to drop
	for i, j := range p.data.Header {
		if p.filter.DropColumn(strings.ToUpper(j.VariableName)) {
			rowsToDrop[i] = true
			continue
		}
		headers = append(headers, j)
		rowsToDrop[i] = false
	}

	// now drop them
	for i, j := range p.data.Rows {
		newRowData := make([]string, 0, p.data.HeaderCount)
		for col, z := range j.RowData {
			if rowsToDrop[col] {
				continue
			}
			newRowData = append(newRowData, z)
		}
		p.data.Rows[i].RowData = newRowData
	}

	p.data.Header = headers
	p.data.HeaderCount = len(headers)

	p.audit.NumObLoaded = p.data.RowCount
	p.audit.NumVarLoaded = p.data.HeaderCount

	return nil
}
