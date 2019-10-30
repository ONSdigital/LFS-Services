package mysql

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/config"
	"services/dataset"
	"services/types"
	"services/util"
	"strconv"
	"strings"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var surveyTable string

func init() {
	surveyTable = config.Config.Database.SurveyTable
	if surveyTable == "" {
		panic("survey table configuration not set")
	}
}

func (s MySQL) DeleteSurveyData(name string) (bool, error) {
	//cnt, err := s.DB.Collection(surveyTable).Find("file_name", name).Count()
	cnt, err := s.DB.Collection(surveyTable).Find(db.Cond{"file_name": name}).Count()

	if err != nil {
		return false, err
	}
	if cnt == 0 {
		return false, nil
	}
	q := s.DB.DeleteFrom(surveyTable).Where("file_name", name)
	_, err = q.Exec()
	if err != nil {
		return true, err
	}
	return true, nil
}

func (s MySQL) insertSurveyData(tx sqlbuilder.Tx, survey types.Survey) error {

	col := tx.Collection(surveyTable)
	_, err := col.Insert(survey)
	if err != nil {
		return err
	}

	return nil
}

func (s MySQL) UnpersistSurveyDataset(tableName string) (dataset.Dataset, error) {
	d, err := dataset.NewDataset(tableName)

	startTime := time.Now()
	log.Info().Msg("starting unpersist")

	if err != nil {
		log.Error().
			Err(err).
			Str("methodName", "UnpersistSurveyDataset").
			Msg("Cannot create a new DataSet")
		return dataset.Dataset{}, fmt.Errorf("cannot create a new DataSet: %s", err)
	}

	log.Info().Msg("starting unpersist into Dataset")

	req := s.DB.Collection(surveyTable).Find().
		Where("file_name = '" + tableName + "'").
		OrderBy("column_number")

	var column types.Survey
	for req.Next(&column) {
		a := strings.Split(column.Rows, ",")
		s := make([]interface{}, len(a))
		for i, v := range a {
			switch reflect.Kind(column.Kind) {
			case reflect.String:
				s[i] = v
			case reflect.Int8, reflect.Uint8:
				s[i], err = strconv.ParseInt(v, 10, 8)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int8 - possible corruption")
				}
			case reflect.Int, reflect.Int32, reflect.Uint32:
				s[i], err = strconv.ParseInt(v, 10, 32)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int32 - possible corruption")
				}
			case reflect.Int64, reflect.Uint64:
				s[i], err = strconv.ParseInt(v, 10, 64)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on int64 - possible corruption")
				}
			case reflect.Float32:
				s[i], err = strconv.ParseFloat(v, 32)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on float32 - possible corruption")
				}
			case reflect.Float64:
				s[i], err = strconv.ParseFloat(v, 64)
				if err != nil {
					return dataset.Dataset{}, fmt.Errorf(" -> UnpersistSurveyDataset: unpersist error on float64 - possible corruption")
				}
			default:
				log.Error().
					Err(err).
					Str("methodName", "UnpersistSurveyDataset").
					Str("type", string(reflect.Kind(column.Kind))).
					Msg("Unknown type - possible corruption")
				return dataset.Dataset{}, fmt.Errorf("unknown type - possible corruption")
			}
		}
		col := dataset.Column{
			ColNo: column.ColumnNumber,
			Kind:  reflect.Kind(column.Kind),
			Rows:  s,
		}
		d.Columns[column.ColumnName] = col
		d.ColumnCount++
		d.RowCount = len(s)
	}

	log.Debug().
		Err(err).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Data unpersisted")
	return d, nil
}

func (s MySQL) PersistSurveyDataset(d dataset.Dataset, vo types.SurveyVO) error {
	var kBuffer bytes.Buffer

	startTime := time.Now()
	log.Debug().Msg("Starting persistence into DB")

	found, err := s.DeleteSurveyData(vo.FileName)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Delete existing survey data failed")
		return fmt.Errorf("delete existing survey data failed, error: %s", err)
	}

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	for colName, column := range d.Columns {
		kBuffer.Reset()

		var i = 0
		columnKind := column.Kind
		for _, v := range column.Rows {
			switch columnKind {
			case reflect.String:
				a := fmt.Sprintf("%s", v)
				a = strings.Replace(a, "'", `''`, -1)
				//kBuffer.WriteString("\"" + a + "\"")
				kBuffer.WriteString(a)
			case reflect.Int8, reflect.Uint8:
				kBuffer.WriteString(fmt.Sprintf("%d", v))
			case reflect.Int, reflect.Int32, reflect.Uint32:
				kBuffer.WriteString(fmt.Sprintf("%d", v))
			case reflect.Int64, reflect.Uint64:
				kBuffer.WriteString(fmt.Sprintf("%d", v))
			case reflect.Float32:
				kBuffer.WriteString(fmt.Sprintf("%f", v))
			case reflect.Float64:
				kBuffer.WriteString(fmt.Sprintf("%g", v))
			default:
				log.Error().
					Str("methodName", "PersistSurveyDataset").
					Int("type", int(columnKind)).
					Msg("Unknown type - possible corruption")
				return fmt.Errorf("unknown type - possible corruption")
			}

			if i != len(column.Rows)-1 {
				kBuffer.WriteString(",")
			}
			i++
		}

		row := types.Survey{
			Id:           vo.Id,
			FileName:     vo.FileName,
			FileSource:   vo.FileSource,
			Week:         vo.Week,
			Month:        vo.Month,
			Year:         vo.Year,
			ColumnName:   colName,
			ColumnNumber: column.ColNo,
			Kind:         int(columnKind),
			Rows:         kBuffer.String(),
		}

		if err := s.insertSurveyData(tx, row); err != nil {
			return fmt.Errorf("cannot insert survey row, error: %s", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	vo.NumObLoaded = d.NumObLoaded
	vo.NumObFile = d.NumObFile
	vo.NumVarLoaded = d.NumVarLoaded
	vo.NumVarFile = d.NumVarFile

	weekOrMonth := func() int {
		if vo.FileSource == types.GBSource {
			return vo.Week
		}
		return vo.Month
	}()

	var updateMonthlyStatus func(int, int, int) error
	if vo.FileSource == types.GBSource {
		updateMonthlyStatus = s.updateGBBatch
	} else {
		updateMonthlyStatus = s.updateNIBatch
	}

	if found {
		if err := updateMonthlyStatus(weekOrMonth, vo.Year, types.FileReloaded); err != nil {
			return err
		}
		if err := s.insertAudit(vo, types.FileReloaded, "File re-uploaded successfully"); err != nil {
			return err
		}
	} else {
		if err := updateMonthlyStatus(weekOrMonth, vo.Year, types.FileUploaded); err != nil {
			return err
		}
		if err := s.insertAudit(vo, types.FileUploaded, "File uploaded successfully"); err != nil {
			return err
		}
	}

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("SurveyInput data persisted")

	return nil
}

func (s MySQL) insertAudit(vo types.SurveyVO, status int, message string) error {
	audit := types.Audit{
		Id:            vo.Id,
		FileName:      surveyTable,
		FileSource:    vo.FileSource,
		Week:          vo.Week,
		Month:         vo.Month,
		Year:          vo.Year,
		ReferenceDate: time.Now(),
		NumVarFile:    vo.NumVarFile,
		NumVarLoaded:  vo.NumVarLoaded,
		NumObFile:     vo.NumObFile,
		NumObLoaded:   vo.NumObLoaded,
		Status:        status,
		Message:       message,
	}

	var f = DBAudit{s}
	if err := f.AuditFileUploadEvent(audit); err != nil {
		log.Error().
			Err(err).
			Msg("AuditFileUpload failed")
		return fmt.Errorf("AuditFileUpload, error: %s", err)
	}
	return nil
}
