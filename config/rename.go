package config

type Rename struct {
	Survey []Columns
}

type Columns struct {
	From string
	To   string
}
