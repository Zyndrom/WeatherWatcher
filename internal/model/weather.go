package model

type Weather struct {
	Location      Location
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	WindSpeed     float32 `json:"wind_speed"`
	Temp          float32 `json:"temp"`
	TempFeelsLike float32 `json:"temp_feels_like"`
	Humidity      float32 `json:"humidity"`
}
