package service

import (
	"weather-app/entity"
	utilapi "weather-app/util-api"
)

type WeatherService interface {
	GetWeather(string) (entity.Weather, error, bool, any)
}

type weatherService struct {
}

func New() WeatherService {
	return &weatherService{}
}

func (service *weatherService) GetWeather(location string) (entity.Weather, error, bool, any) {
	responseData, err, isInvalidResponse, invalidResponse := utilapi.GetWeather(location)

	if err != nil {
		return entity.Weather{}, err, false, nil
	} else if isInvalidResponse {
		return entity.Weather{}, nil, true, invalidResponse
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
		return weatherData, nil, false, nil
	}
}
