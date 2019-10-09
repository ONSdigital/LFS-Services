package config

type DropColumns struct {
	Survey DropColumnNames
}

type DropColumnNames struct {
	ColumnNames []string
}

// [dropColumns.survey] columnNames = ["SOC2KA_INDEX", "HLDCOUNT"]
