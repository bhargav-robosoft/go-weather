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

func (controller *controller) GetWeather(ctx *gin.Context) (weather entity.Weather, err error) {
	params := ctx.Request.URL.Query()
	if !params.Has("location") {
		return weather, errors.New("No location")
	}
	location := params["location"][0]

	return controller.service.GetWeather(location), nil
}
