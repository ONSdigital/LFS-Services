package config

type DropColumns struct {
	Survey DropColumnNames
}

type DropColumnNames struct {
	ColumnNames []string
}
