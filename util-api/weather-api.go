package utilapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

func GetWeather(location string) WeatherApiResponseWeatherData {
	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?&appid=e07a248e19b7bc76072304519cc9e7ff&units=metric&q=" + location)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject WeatherApiResponseWeatherData
	njson.Unmarshal(responseData, &responseObject)

	return responseObject
}
