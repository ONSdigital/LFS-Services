package util

import "database/sql"

//ToNullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
