package utilapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/m7shapan/njson"
)

type WeatherApiResponseWeatherData struct {
	LocationName        string  `njson:"name"`
	LocationCountryName string  `njson:"sys.country"`
	Temperature         float64 `njson:"main.temp"`
	Description         string  `njson:"weather.0.description"`
	WeatherIcon         string  `njson:"weather.0.icon"`
	MinTemperature      float64 `njson:"main.temp_min"`
	MaxTemperature      float64 `njson:"main.temp_max"`
	Clouds              int     `njson:"clouds.all"`
	Humidity            int     `njson:"main.humidity"`
	WindSpeed           float64 `njson:"wind.speed"`
	Visibility          int     `njson:"visibility"`
}

func GetWeather(location string) (WeatherApiResponseWeatherData, error) {
	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?&appid=e07a248e19b7bc76072304519cc9e7ff&units=metric&q=" + location)

	if err != nil {
		return WeatherApiResponseWeatherData{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return WeatherApiResponseWeatherData{}, err
	}

	if response.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)
		return WeatherApiResponseWeatherData{}, errors.New(errorResponse["message"])
	}

	var responseObject WeatherApiResponseWeatherData
	err = njson.Unmarshal(responseData, &responseObject)
	if err != nil {
		fmt.Println("error", err)
		return WeatherApiResponseWeatherData{}, err
	}
	fmt.Println("response", responseObject)

	return responseObject, nil
}
