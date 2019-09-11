package spss

type ColumnType int32
type ColumnTypes string

const (
	ReadstatTypeString    ColumnType = 0
	ReadstatTypeInt8      ColumnType = 1
	ReadstatTypeInt16     ColumnType = 2
	ReadstatTypeInt32     ColumnType = 3
	ReadstatTypeFloat     ColumnType = 4
	ReadstatTypeDouble    ColumnType = 5
	ReadstatTypeStringRef ColumnType = 6
)

const (
	INT     ColumnTypes = "INTEGER"
	INTEGER ColumnTypes = "INTEGER"
	BIGINT  ColumnTypes = "BIGINT"
	STRING  ColumnTypes = "TEXT"
	FLOAT   ColumnTypes = "REAL"
	DOUBLE  ColumnTypes = "REAL"
)

func (columnType ColumnType) As() ColumnType {
	return columnType
}

func (columnType ColumnType) AsInt() int {
	return int(columnType)
}

func (columnType ColumnType) IsNumeric() bool {
	if columnType != ReadstatTypeInt8 && columnType != ReadstatTypeInt16 &&
		columnType != ReadstatTypeInt32 && columnType != ReadstatTypeFloat &&
		columnType != ReadstatTypeDouble {
		return false
	}
	return true
}
