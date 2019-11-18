package postgres

import (
	"github.com/rs/zerolog/log"
)

func (s Postgres) DeleteFrom(table string) error {
	q := s.DB.DeleteFrom(table)

	_, err := q.Exec()

	if err != nil {
		log.Debug().
			Str("table", table).
			Msg("Delete from table failed: " + err.Error())
		return err
	}

	return err
}

func (s Postgres) CountRows(table string) (int, error) {
	var count int

	row, _ := s.DB.QueryRow("SELECT COUNT(*) FROM " + table)
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}
