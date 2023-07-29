package controller

import (
	"errors"
	"weather-app/entity"
	"weather-app/service"

	"github.com/gin-gonic/gin"
)

type WeatherController interface {
	GetWeather(ctx *gin.Context) (entity.Weather, error, bool, any)
}

type controller struct {
	service service.WeatherService
}

func New(service service.WeatherService) WeatherController {
	return &controller{
		service: service,
	}
}

func (controller *controller) GetWeather(ctx *gin.Context) (entity.Weather, error, bool, any) {
	params := ctx.Request.URL.Query()

	if !params.Has("location") {
		return entity.Weather{}, errors.New("No location"), false, nil
	}

	location := params["location"][0]

	if len(location) == 0 {
		return entity.Weather{}, errors.New("Empty location"), false, nil
	}

	data, err, isInvalidResponse, invalidResponse := controller.service.GetWeather(location)
	if err != nil {
		return entity.Weather{}, err, false, nil
	} else if isInvalidResponse {
		return entity.Weather{}, err, true, invalidResponse
	}

	return data, nil, false, nil
}
