package types

type Columns struct {
	TableName    string `db:"table_name"`
	ColumnName   string `db:"column_name"`
	ColumnNumber int    `db:"column_number"`
	Kind         int    `db:"kind"`
	Rows         string `db:"rows"`
}
