package api

import (
	"net/http"
	"weather-app/controller"
	"weather-app/entity"
	"weather-app/service"

	"github.com/gin-gonic/gin"
)

var (
	weatherService    service.WeatherService       = service.New()
	weatherController controller.WeatherController = controller.New(weatherService)
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gin.Default()

	server.GET("/get-weather", func(ctx *gin.Context) {
		data, err := weatherController.GetWeather(ctx)
		if err != nil {
			errorData := entity.WeatherError{
				Status:  404,
				Message: "Location not set",
				Example: "/get-weather?location=udupi",
			}
			ctx.JSON(404, errorData)
		} else {
			ctx.JSON(200, data)
		}
	})

	server.ServeHTTP(w, r)
}
