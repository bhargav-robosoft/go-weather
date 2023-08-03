package service

import (
	"errors"
	"fmt"
	"weather-app/db"
	"weather-app/entity"
	utils "weather-app/utils"
)

type WeatherService interface {
	getWeather(location string) (entity.Weather, error)
	AddToRecent(location string, cookieId string) (weatherData entity.Weather, id string, err error)
	GetRecentsWeather(cookieId string) ([]entity.Weather, error)
	GetFavouritesWeather(cookieId string) ([]entity.Weather, error)
}

type weatherService struct{}

func New() WeatherService {
	return &weatherService{}
}

func (service *weatherService) getWeather(location string) (entity.Weather, error) {
	if len(location) == 0 {
		return entity.Weather{}, errors.New("Empty location")
	}

	responseData, err := utils.GetWeather(location)
	fmt.Println("responseData", responseData, err)

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

func (service *weatherService) AddToRecent(location string, id string) (entity.Weather, string, error) {
	weatherData, err := service.getWeather(location)
	fmt.Println("Weather1", weatherData, err)
	if err != nil {
		return entity.Weather{}, "", err
	}

	dbId, err := db.AddRecentSearchForUser(id, weatherData.Name)
	fmt.Println("AddRecentSearchForUser", dbId, err)
	if err != nil {
		return weatherData, dbId, err
	}

	if dbId != id {
		return weatherData, dbId, nil
	}

	isFav, err := db.IsFavourite(dbId, weatherData.Name)
	fmt.Println("IsFavourite", isFav, err)
	if err != nil {
		return weatherData, dbId, err
	}

	weatherData.IsFavourite = isFav

	return weatherData, dbId, err
}

func (service *weatherService) GetRecentsWeather(cookieId string) ([]entity.Weather, error) {
	recentLocations, favouriteLocations, err := db.GetRecentsAndFavouritesForUser(cookieId)
	if err != nil {
		return []entity.Weather{}, err
	}

	var recentsWeatherData = []entity.Weather{}

	for _, location := range recentLocations {
		weatherData, err := service.getWeather(location)
		if err != nil {
			continue
		} else {
			if utils.Contains(favouriteLocations, location) {
				weatherData.IsFavourite = true
			}
			recentsWeatherData = append(recentsWeatherData, weatherData)
		}
	}
	return recentsWeatherData, nil
}

func (service *weatherService) GetFavouritesWeather(cookieId string) ([]entity.Weather, error) {
	favouriteLocations, err := db.GetFavouritesForUser(cookieId)
	if err != nil {
		return []entity.Weather{}, err
	}

	var favouritesWeatherData = []entity.Weather{}

	for _, v := range favouriteLocations {
		weatherData, err := service.getWeather(v)
		if err != nil {
			continue
		} else {
			weatherData.IsFavourite = true
			favouritesWeatherData = append(favouritesWeatherData, weatherData)
		}
	}
	return favouritesWeatherData, nil
}
