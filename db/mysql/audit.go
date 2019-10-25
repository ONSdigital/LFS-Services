package mysql

import (
	"services/config"
	"services/dataset"
	"services/types"
	"time"
)

type DBAudit struct {
	MySQL
}

var surveyAuditTable string

func init() {
	surveyAuditTable = config.Config.Database.SurveyAuditTable
	if surveyAuditTable == "" {
		panic("survey audit table configuration not set")
	}
}

func (s MySQL) AuditFileUploadEvent(d dataset.Dataset, id int) error {
	event := types.Audit{
		Id:            id,
		FileName:      d.DatasetName,
		ReferenceDate: time.Now(),
		NumVarFile:    d.NumVarFile,
		NumVarLoaded:  d.NumVarLoaded,
		NumObFile:     d.NumObFile,
		NumObLoaded:   d.NumObLoaded,
	}
	dbAudit := s.DB.Collection(surveyAuditTable)
	_, err := dbAudit.Insert(event)
	if err != nil {
		return err
	}

	return nil
}
