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
			errorData := entity.WeatherFailureResponse{
				Status:  404,
				Message: err.Error(),
				Example: "/get-weather?location=udupi",
			}
			ctx.JSON(404, errorData)
		} else {
			successData := entity.WeatherSuccessResponse{
				Status:  200,
				Message: "Fetched weather data",
				Data:    data,
			}
			ctx.JSON(200, successData)
		}
	})

	server.GET("/get-recents", func(ctx *gin.Context) {
		data, err := weatherController.GetRecentsWeather(ctx)
		if err != nil {
			errorData := entity.MultiWeatherResponse{
				Status:  404,
				Message: err.Error(),
			}
			ctx.JSON(404, errorData)
		} else {
			successData := entity.MultiWeatherResponse{
				Status:  200,
				Message: "Fetched recents",
				Data:    data,
			}
			ctx.JSON(200, successData)
		}
	})

	server.GET("/get-favourites", func(ctx *gin.Context) {
		data, err := weatherController.GetFavouritesWeather(ctx)
		if err != nil {
			errorData := entity.MultiWeatherResponse{
				Status:  404,
				Message: err.Error(),
			}
			ctx.JSON(404, errorData)
		} else {
			successData := entity.MultiWeatherResponse{
				Status:  200,
				Message: "Fetched favourites",
				Data:    data,
			}
			ctx.JSON(200, successData)
		}
	})

	server.GET("/handle-favourite", func(ctx *gin.Context) {
		response, err := weatherController.HandleFavourite(ctx)
		if err != nil {
			errorData := entity.NormalResponse{
				Status:  404,
				Message: err.Error(),
			}
			ctx.JSON(404, errorData)
		} else {
			successData := entity.NormalResponse{
				Status:  200,
				Message: response,
			}
			ctx.JSON(200, successData)
		}
	})

	server.GET("/clear-recents", func(ctx *gin.Context) {
		response, err := weatherController.ClearRecents(ctx)
		if err != nil {
			errorData := entity.NormalResponse{
				Status:  404,
				Message: err.Error(),
			}
			ctx.JSON(404, errorData)
		} else {
			data := entity.NormalResponse{
				Status:  200,
				Message: response,
			}
			ctx.JSON(200, data)
		}
	})

	server.GET("/clear-favourites", func(ctx *gin.Context) {
		response, err := weatherController.ClearFavourites(ctx)
		if err != nil {
			errorData := entity.NormalResponse{
				Status:  404,
				Message: err.Error(),
			}
			ctx.JSON(404, errorData)
		} else {
			data := entity.NormalResponse{
				Status:  200,
				Message: response,
			}
			ctx.JSON(200, data)
		}
	})

	server.Run(":8080")
}
