package service

import (
	"weather-app/entity"
	utilapi "weather-app/util-api"
)

type WeatherService interface {
	GetWeather(string) (entity.Weather, error)
	GetRecents([]string) ([]entity.Weather, error)
	GetFavourites([]string) ([]entity.Weather, error)
}

type weatherService struct{}

func New() WeatherService {
	return &weatherService{}
}

func (service *weatherService) GetWeather(location string) (entity.Weather, error) {
	responseData, err := utilapi.GetWeather(location)

	if err != nil {
		return entity.Weather{}, err
	} else {
		var weatherData entity.Weather
		weatherData = entity.Weather{
			Name:            responseData.LocationName,
			CountryName:     responseData.LocationCountryName,
			Temperature:     responseData.Temperature,
			Description:     responseData.Description,
			WeatherIconLink: "https://openweathermap.org/img/wn/" + responseData.WeatherIcon + "@2x.png",
			MinTemperature:  responseData.MinTemperature,
			MaxTemperature:  responseData.MaxTemperature,
			Clouds:          responseData.Clouds,
			Humidity:        responseData.Humidity,
			WindSpeed:       responseData.WindSpeed,
			Visibility:      responseData.Visibility,
		}
		return weatherData, nil
	}
}

func (service *weatherService) GetRecents(locations []string) ([]entity.Weather, error) {
	var recentWeatherData = []entity.Weather{}

	for _, v := range locations {
		weatherData, err := service.GetWeather(v)
		if err != nil {
			continue
		} else {
			recentWeatherData = append(recentWeatherData, weatherData)
		}

	}
	return recentWeatherData, nil
}

func (service *weatherService) GetFavourites(locations []string) ([]entity.Weather, error) {
	var favouritesWeatherData = []entity.Weather{}

	for _, v := range locations {
		weatherData, err := service.GetWeather(v)
		if err != nil {
			continue
		} else {
			favouritesWeatherData = append(favouritesWeatherData, weatherData)
		}

	}
	return favouritesWeatherData, nil
}
