package model

type Location struct {
	Id      uint32  `json:"id"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Name    string  `json:"name"`
	RuName  string  `json:"ru_name"`
	Country string  `json:"country"`
}
