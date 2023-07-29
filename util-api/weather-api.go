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
	LocationName        string `njson:"name"`
	LocationCountryName string `njson:"sys.country"`
	Temperature         string `njson:"main.temp"`
	Description         string `njson:"weather.0.description"`
	WeatherIcon         string `njson:"weather.0.icon"`
	MinTemperature      string `njson:"main.temp_min"`
	MaxTemperature      string `njson:"main.temp_max"`
	Clouds              string `njson:"clouds.all"`
	Humidity            string `njson:"main.humidity"`
	WindSpeed           string `njson:"wind.speed"`
	Visibility          string `njson:"visibility"`
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
