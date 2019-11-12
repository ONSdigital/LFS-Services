package mysql

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"services/config"
	"services/types"
	"services/util"
	"time"
)

var addressesTable string

const BatchSize = 5000

func init() {
	addressesTable = config.Config.Database.AddressesTable
	if addressesTable == "" {
		panic("addresses table configuration not set")
	}
}

func (s MySQL) DeleteAddressesData(name string) error {
	return s.DB.Collection(addressesTable).Truncate()
}

func (s MySQL) insertAddressesRow(buffer bytes.Buffer) error {

	tx, err := s.DB.NewTx(nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Start transaction failed")
		return fmt.Errorf("cannot start a transaction, error: %s", err)
	}

	_, err = s.DB.Exec(buffer.String())

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().
			Err(err).
			Msg("Commit transaction failed")
		return fmt.Errorf("commit failed, error: %s", err)
	}

	return nil
}

func (s MySQL) PersistAddressDataset(header []string, rows [][]string, status *types.WSMessage) error {

	startTime := time.Now()

	log.Debug().
		Str("tableName", addressesTable).
		Msg("Starting persistence into DB")

	_ = s.DeleteAddressesData(addressesTable)

	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO " + addressesTable + "(")
	for i := 0; i < len(header); i++ {
		buffer.WriteString(header[i])
		if i != len(header)-1 {
			buffer.WriteString(",")
		} else {
			buffer.WriteString(") VALUES ")
		}
	}

	t1 := reflect.TypeOf(types.Addresses{})

	meta := make([]reflect.Kind, t1.NumField())
	for i := 0; i < t1.NumField(); i++ {
		a := t1.Field(i)
		meta[i] = a.Type.Kind()
	}

	cnt := 0
	batchCount := 0

	for j, row := range rows {
		buffer.WriteString("(")

		for i, v := range row {
			columnKind := meta[i]
			if v == "" {
				switch columnKind {
				case reflect.String:
					buffer.WriteString("NULL")
				case reflect.Int8, reflect.Uint8:
					buffer.WriteString("0")
				case reflect.Int, reflect.Int32, reflect.Uint32:
					buffer.WriteString("0")
				case reflect.Int64, reflect.Uint64:
					buffer.WriteString("0")
				case reflect.Float32:
					buffer.WriteString("0.0")
				case reflect.Float64:
					buffer.WriteString("0.0")
				default:
					log.Error().
						Str("methodName", "PersistAddresses").
						Int("type", int(columnKind)).
						Msg("Unknown type - possible corruption")
					return fmt.Errorf("unknown type - possible corruption")
				}
			} else {
				if columnKind == reflect.String {
					buffer.WriteString("'" + v + "'")
				} else {
					buffer.WriteString(v)
				}
			}
			if i != len(header)-1 {
				buffer.WriteString(",")
			} else {
				buffer.WriteString(") ")
			}
		}

		cnt++
		if cnt == BatchSize {
			err := s.insertAddressesRow(buffer)
			if err != nil {
				log.Error().
					Err(err).
					Msg("insert addreses failed")
				return fmt.Errorf("cannot insert an addresses record, error: %s", err)
			}
			cnt = 0
			batchCount = batchCount + 1
			buffer.Reset()
			buffer.WriteString("INSERT INTO " + addressesTable + "(")
			for i := 0; i < len(header); i++ {
				buffer.WriteString(header[i])
				if i != len(header)-1 {
					buffer.WriteString(",")
				} else {
					buffer.WriteString(") VALUES ")
				}
			}

			var perc = (float64(batchCount*BatchSize) / float64(len(rows))) * 100
			status.SetPercentage(perc)
		} else {
			if j != len(rows)-1 {
				buffer.WriteString(",")
			}
		}
	}

	if cnt > 0 {
		err := s.insertAddressesRow(buffer)
		if err != nil {
			log.Error().
				Err(err).
				Msg("insert addreses failed")
			status.SetUploadError(fmt.Sprintf("cannot insert an addresses record, error: %s", err))
			return fmt.Errorf("cannot insert an addresses record, error: %s", err)
		}
	}

	log.Debug().
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("Addresses data persisted")

	status.SetUploadFinished()

	return nil
}
