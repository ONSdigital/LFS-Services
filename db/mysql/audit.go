package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
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

func (s MySQL) AuditFileUploadEvent(event types.Audit) error {

	dbAudit := s.DB.Collection(surveyAuditTable)
	_, err := dbAudit.Insert(event)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) GetAllAudits() ([]types.Audit, error) {

	var audits []types.Audit
	res := s.DB.Collection(surveyAuditTable).Find()
	defer func() { _ = res.Close() }()
	err := res.All(&audits)
	if err != nil {
		return nil, res.Err()
	}

	return audits, nil
}

func (s MySQL) GetAuditsByYear(year types.Year) ([]types.Audit, error) {

	var audits []types.Audit
	dbAudit := s.DB.Collection(surveyAuditTable)
	if dbAudit.Exists() {
		log.Error().Str("table", surveyAuditTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", surveyAuditTable)
	}
	res := dbAudit.Find("year", year)
	defer func() { _ = res.Close() }()
	if res.Err() != nil {
		return nil, res.Err()
	}

	return audits, nil
}

func (s MySQL) GetAuditsByYearMonth(month types.Month, year types.Year) ([]types.Audit, error) {

	var audits []types.Audit
	dbAudit := s.DB.Collection(surveyAuditTable)
	if dbAudit.Exists() {
		log.Error().Str("table", surveyAuditTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", surveyAuditTable)
	}
	res := dbAudit.Find("year", year, "month", month)
	defer func() { _ = res.Close() }()
	if res.Err() != nil {
		return nil, res.Err()
	}

	return audits, nil
}

func (s MySQL) GetAuditsByYearWeek(week types.Week, year types.Year) ([]types.Audit, error) {

	var audits []types.Audit
	dbAudit := s.DB.Collection(surveyAuditTable)
	if dbAudit.Exists() {
		log.Error().Str("table", surveyAuditTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", surveyAuditTable)
	}
	res := dbAudit.Find("year", year, "week", week)
	defer func() { _ = res.Close() }()
	if res.Err() != nil {
		return nil, res.Err()
	}

	return audits, nil
}
