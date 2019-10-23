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

var uploadAuditTable string

func init() {
	uploadAuditTable = config.Config.Database.UploadAuditTable
	if uploadAuditTable == "" {
		panic("upload audit table configuration not set")
	}
}

func (s MySQL) AuditFileUploadEvent(d dataset.Dataset) error {
	event := types.Audit{
		FileName:      d.DatasetName,
		ReferenceDate: time.Now(),
		NumVarFile:    d.NumVarFile,
		NumVarLoaded:  d.NumVarLoaded,
		NumObFile:     d.NumObFile,
		NumObLoaded:   d.NumObLoaded,
	}
	dbAudit := s.DB.Collection(uploadAuditTable)
	_, err := dbAudit.Insert(event)
	if err != nil {
		return err
	}

	return nil
}
