package types

type Dashboard struct {
	Id     int    `db:"id" json:"id"`
	Type   string `json:"type"`
	Period string `db:"quarter" db:"month" json:"period"`
	Year   int    `db:"year" json:"year"`
	Status int    `db:"status" json:"status"`
}
