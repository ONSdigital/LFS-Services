package types

type MonthlyBatch struct {
	Id          int    `db:"id,omitempty"`
	Year        int    `db:"year"`
	Month       int    `db:"month"`
	Status      int    `db:"status"`
	Description string `db:"description"`
}

type GBBatchItem struct {
	Id     int `db:"id"`
	Year   int `db:"year"`
	Month  int `db:"month"`
	Week   int `db:"week"`
	Status int `db:"status"`
}

type NIBatchItem struct {
	Id     int `db:"id"`
	Year   int `db:"year"`
	Month  int `db:"month"`
	Status int `db:"status"`
}

type QuarterlyBatch struct {
	Id          int    `db:"id,omitempty"`
	Quarter     int    `db:"quarter"`
	Year        int    `db:"year"`
	Status      int    `db:"status"`
	Description string `db:"description"`
}

type AnnualBatch struct {
	Id          int    `db:"id,omitempty"`
	Year        int    `db:"year"`
	Status      int    `db:"status"`
	Description string `db:"description"`
}
