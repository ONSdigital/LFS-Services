package types

// links to the batch_info view
type BatchInfo struct {
	Id          int    `db:"id"`
	MMonth      int    `db:"m_month"`
	MYear       int    `db:"m_year"`
	MStatus     int    `db:"m_status"`
	Description string `db:"description"`
	NIYear      int    `db:"ni_year"`
	NIMonth     int    `db:"ni_month"`
	NIStatus    int    `db:"ni_status"`
	GBYear      int    `db:"gb_year"`
	GBMonth     int    `db:"gb_month"`
	GBWeek      int    `db:"gb_month"`
	GBStatus    int    `db:"gb_status"`
}

type GBBatchInfo struct {
	Id          int    `db:"id"`
	Year        int    `db:"year"`
	Month       int    `db:"month"`
	Status      int    `db:"status"`
	Description string `db:"description"`
	Week        int    `db:"week"`
}

type NIBatchInfo struct {
	Id          int    `db:"id"`
	Year        int    `db:"year"`
	Month       int    `db:"month"`
	Status      int    `db:"status"`
	Description string `db:"description"`
}

type MonthlyBatch struct {
	Id          int    `db:"id"`
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
