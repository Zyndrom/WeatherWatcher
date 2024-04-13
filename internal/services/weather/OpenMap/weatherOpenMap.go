package openMapWeather

import (
	"GoWeatherMap/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type OpenMapWeather struct {
	apiKey string
}

func New(apiKey string) *OpenMapWeather {
	return &OpenMapWeather{apiKey: apiKey}
}

func (w *OpenMapWeather) FindLocationByName(name string) ([]model.Location, error) {
	query := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&appid=%s&units=metric&limit=5", name, w.apiKey)
	resp, err := http.Get(query)
	if err != nil || resp == nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		logrus.Warningf("Can not unmarshal JSON, %s", string(body))
		return nil, err
	}
	type Anss struct {
		Name        string `json:"name"`
		Local_names struct {
			Ru string `json:"ru"`
		} `json:"local_names"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		CountryCode string  `json:"country"`
	}
	var ans []Anss
	if err := json.Unmarshal(body, &ans); err != nil {
		logrus.Warningf("Can not unmarshal JSON, %s", string(body))
		return nil, err
	}
	res := []model.Location{}
	for _, val := range ans {
		loc := model.Location{
			Lat:     val.Lat,
			Lon:     val.Lon,
			Name:    val.Name,
			RuName:  val.Local_names.Ru,
			Country: val.CountryCode,
		}
		res = append(res, loc)
	}
	return res, nil
}

func (w *OpenMapWeather) GetWeather(loc model.Location) (model.Weather, error) {
	fmt.Println(loc)
	query := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", loc.Lat, loc.Lon, w.apiKey)
	resp, err := http.Get(query)
	weather := model.Weather{}
	if err != nil || resp == nil {
		log.Println(err)
		return weather, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		logrus.Warningf("Can not unmarshal JSON, %s", string(body))
		return weather, err
	}
	var ans struct {
		Weather []struct {
			Main string `json:"main"`
			Desc string `json:"description"`
		} `json:"weather"`
		Main struct {
			Temp      float32 `json:"temp"`
			FeelsLike float32 `json:"feels_like"`
			Humidity  float32 `json:"Humidity"`
		} `json:"main"`
		Wind struct {
			Speed float32 `json:"speed"`
		} `json:"Wind"`
	}
	if err := json.Unmarshal(body, &ans); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return weather, err
	}
	if len(ans.Weather) < 1 {
		return weather, err
	}
	weather.Location = loc
	weather.Name = ans.Weather[0].Main
	weather.Description = ans.Weather[0].Desc
	weather.WindSpeed = ans.Wind.Speed
	weather.Temp = ans.Main.Temp
	weather.TempFeelsLike = ans.Main.FeelsLike
	weather.Humidity = ans.Main.Humidity

	return weather, nil
}
