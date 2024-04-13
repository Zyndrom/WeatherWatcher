package weather

import (
	"GoWeatherMap/internal/config"
	"GoWeatherMap/internal/model"
	openMapWeather "GoWeatherMap/internal/services/weather/OpenMap"
	"strings"
)

type WeatherService struct {
	weather *openMapWeather.OpenMapWeather
	storage storage
}
type storage interface {
	GetLocationByName(name string) ([]model.Location, error)
	SetLocationByName(name string, locs []model.Location) error
}

func NewWeatherController(cfg *config.Config, storage storage) *WeatherService {
	return &WeatherService{weather: openMapWeather.New(cfg.OpenMapKey), storage: storage}
}

func (w *WeatherService) FindLocationByName(name string) ([]model.Location, error) {
	name = strings.ReplaceAll(name, " ", "-")
	locs, err := w.storage.GetLocationByName(name)
	if err != nil || len(locs) == 0 {
		locs, err = w.weather.FindLocationByName(name)
		if err != nil || len(locs) == 0 {
			w.storage.SetLocationByName(name, locs)
		}
	}
	return locs, err
}

func (w *WeatherService) CurrentWeatherInLocation(loc model.Location) (model.Weather, error) {
	return w.weather.GetWeather(loc)
}
