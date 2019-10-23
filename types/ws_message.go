package types

type WSMessage struct {
	Filename   string  `json:"fileName"`
	Percentage float64 `json:"percent"`
	Status     int     `json:"status"`
}
