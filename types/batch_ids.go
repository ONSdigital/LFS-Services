package types

type YearID struct {
	Id          int    `db:"id" json:"id"`
	Year        int    `db:"year" json:"year"`
	Status      int    `db:"status" json:"status"`
	Description string `db:"description" json:"description"`
}

type QuarterID struct {
	Id          int    `db:"id" json:"id"`
	Quarter     int    `db:"quarter" json:"quarter"`
	Year        int    `db:"year" json:"year"`
	Status      int    `db:"status" json:"status"`
	Description string `db:"description" json:"description"`
}
