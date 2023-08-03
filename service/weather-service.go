package service

import (
	"errors"
	"weather-app/db"
	"weather-app/entity"
	utils "weather-app/utils"
)

type WeatherService interface {
	getWeather(location string) (entity.Weather, error)
	AddToRecent(location string, cookieId string) (id string, weatherData entity.Weather, err error)
	GetRecentsWeather(cookieId string) (recentsWeather []entity.Weather, err error)
	GetFavouritesWeather(cookieId string) (favouritesWeather []entity.Weather, err error)
	ClearRecents(cookieId string) (err error)
	ClearFavourites(cookieId string) (err error)
	HandleFavourite(location string, cookieId string) (id string, response string, err error)
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

	if err != nil {
		return entity.Weather{}, err
	} else {
		weatherData := entity.Weather{
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

func (service *weatherService) AddToRecent(location string, id string) (string, entity.Weather, error) {
	weatherData, err := service.getWeather(location)
	if err != nil {
		return "", entity.Weather{}, err
	}

	dbId, err := db.AddRecentSearchForUser(id, weatherData.Name)
	if err != nil {
		return dbId, weatherData, err
	}

	if dbId != id {
		return dbId, weatherData, nil
	}

	isFav, err := db.IsFavourite(dbId, weatherData.Name)
	if err != nil {
		return dbId, weatherData, err
	}

	weatherData.IsFavourite = isFav

	return dbId, weatherData, err
}

func (service *weatherService) HandleFavourite(cookieId string, location string) (id string, response string, err error) {
	weatherData, err := service.getWeather(location)
	if err != nil {
		return "", "", err
	}

	return db.HandleFavouriteForUser(cookieId, weatherData.Name)
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

func (service *weatherService) ClearRecents(cookieId string) (err error) {
	return db.ClearRecents(cookieId)
}

func (service *weatherService) ClearFavourites(cookieId string) (err error) {
	return db.ClearFavourites(cookieId)
}
