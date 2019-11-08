package postgres

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"services/config"
	"services/types"
	"upper.io/db.v3"
)

type DBAudit struct {
	Postgres
}

var surveyAuditTable string

func init() {
	surveyAuditTable = config.Config.Database.SurveyAuditTable
	if surveyAuditTable == "" {
		panic("survey audit table configuration not set")
	}
}

func (s Postgres) AuditFileUploadEvent(event types.Audit) error {

	dbAudit := s.DB.Collection(surveyAuditTable)
	_, err := dbAudit.Insert(event)
	if err != nil {
		return err
	}

	return nil
}

func (s Postgres) GetAllAudits() ([]types.Audit, error) {

	var audits []types.Audit
	res := s.DB.Collection(surveyAuditTable).Find()
	defer func() { _ = res.Close() }()
	err := res.All(&audits)
	if err != nil {
		return nil, res.Err()
	}

	return audits, nil
}

func (s Postgres) GetAuditsByYear(year types.Year) ([]types.Audit, error) {

	var audits []types.Audit
	dbAudit := s.DB.Collection(surveyAuditTable)
	if !dbAudit.Exists() {
		log.Error().Str("table", surveyAuditTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", surveyAuditTable)
	}
	//res := dbAudit.Find("year", year)
	res := dbAudit.Find(db.Cond{"year": year})

	defer func() { _ = res.Close() }()
	if res.Err() != nil {
		return nil, res.Err()
	}

	return audits, nil
}

func (s Postgres) GetAuditsByYearMonth(month types.Month, year types.Year) ([]types.Audit, error) {

	var audits []types.Audit
	dbAudit := s.DB.Collection(surveyAuditTable)
	if !dbAudit.Exists() {
		log.Error().Str("table", surveyAuditTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", surveyAuditTable)
	}
	//res := dbAudit.Find("year", year, "month", month)
	res := dbAudit.Find(db.Cond{"year": year, "month": month})
	defer func() { _ = res.Close() }()
	if res.Err() != nil {
		return nil, res.Err()
	}

	return audits, nil
}

func (s Postgres) GetAuditsByYearWeek(week types.Week, year types.Year) ([]types.Audit, error) {

	var audits []types.Audit
	dbAudit := s.DB.Collection(surveyAuditTable)
	if !dbAudit.Exists() {
		log.Error().Str("table", surveyAuditTable).Msg("Table does not exist")
		return nil, fmt.Errorf("table: %s does not exist", surveyAuditTable)
	}
	//res := dbAudit.Find("year", year, "week", week)
	res := dbAudit.Find(db.Cond{"year": year, "week": week})

	defer func() { _ = res.Close() }()
	if res.Err() != nil {
		return nil, res.Err()
	}

	return audits, nil
}
