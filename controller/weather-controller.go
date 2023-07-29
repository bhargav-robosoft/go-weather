package controller

import (
	"errors"
	"weather-app/entity"
	"weather-app/service"

	"github.com/gin-gonic/gin"
)

type WeatherController interface {
	GetWeather(ctx *gin.Context) (entity.Weather, error)
}

type controller struct {
	service service.WeatherService
}

func New(service service.WeatherService) WeatherController {
	return &controller{
		service: service,
	}
}

func (controller *controller) GetWeather(ctx *gin.Context) (entity.Weather, error) {
	params := ctx.Request.URL.Query()

	if !params.Has("location") {
		return entity.Weather{}, errors.New("No location")
	}

	location := params["location"][0]

	if len(location) == 0 {
		return entity.Weather{}, errors.New("Empty location")
	}

	data, err := controller.service.GetWeather(location)
	if err != nil {
		return entity.Weather{}, err
	}

	return data, nil
}
