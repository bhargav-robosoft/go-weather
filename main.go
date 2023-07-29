package main

import (
	"weather-app/controller"
	"weather-app/entity"
	"weather-app/service"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	var (
		weatherService    service.WeatherService       = service.New()
		weatherController controller.WeatherController = controller.New(weatherService)
	)

	server.GET("/get-weather", func(ctx *gin.Context) {
		data, err := weatherController.GetWeather(ctx)
		if err != nil {
			errorData := entity.WeatherError{
				Status:  404,
				Message: err.Error(),
				Example: "/get-weather?location=udupi",
			}
			ctx.JSON(404, errorData)
		} else {
			ctx.JSON(200, data)
		}
	})

	server.Run(":8080")
}
