package types

type MonthlyBatch struct {
	//Id          int    `db:"id"`
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

type AnnualBatch struct {
	Id          int    `db:"id"`
	Year        int    `db:"year"`
	Status      int    `db:"status"`
	Description string `db:"description"`
}
