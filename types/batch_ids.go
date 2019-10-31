package types

type YearID struct {
	Id          int    `db:"id" json:"id"`
	Year        int    `db:"year" json:"year"`
	Status      int    `db:"status" json:"status"`
	Description string `db:"description" json:"description"`
}
