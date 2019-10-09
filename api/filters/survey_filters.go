package filters

import (
	log "github.com/sirupsen/logrus"
	conf "services/config"
	"services/dataset"
)

func NewSurveyFilter(log *log.Logger) *SurveyFilter {
	return &SurveyFilter{log: log}
}

type SurveyFilter struct {
	log *log.Logger
}

func (sf SurveyFilter) DropColumns(dataset *dataset.Dataset) {
	drop := conf.Config.DropColumns.Survey

	for _, v := range drop.ColumnNames {
		sf.log.WithFields(log.Fields{
			"column": v,
		}).Info("Drop column")
		if err := dataset.DropColumn(v); err != nil {
			sf.log.WithFields(log.Fields{
				"column": v,
			}).Warn("drop column not found - ignored")
		}
	}
}

func (sf SurveyFilter) RenameColumns(dataset *dataset.Dataset) {
	cols := conf.Config.Rename.Survey

	for _, v := range cols {
		sf.log.WithFields(log.Fields{
			"from": v.From,
			"to":   v.To,
		}).Info("Rename column")
		if err := dataset.RenameColumn(v.From, v.To); err != nil {
			sf.log.WithFields(log.Fields{
				"from": v.From,
				"to":   v.To,
			}).Warn("Rename column not found - ignored")
		}
	}
}

func (sf SurveyFilter) Validate() error {
	return nil
}
