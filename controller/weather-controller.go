package controller

import (
	"errors"
	"weather-app/entity"
	"weather-app/service"

	"github.com/gin-gonic/gin"
)

type WeatherController interface {
	GetWeather(ctx *gin.Context) (entity.Weather, error)
	GetRecentsWeather(ctx *gin.Context) ([]entity.Weather, error)
	GetFavouritesWeather(ctx *gin.Context) ([]entity.Weather, error)
	HandleFavourite(ctx *gin.Context) (string, error)
	ClearRecents(ctx *gin.Context) (string, error)
	ClearFavourites(ctx *gin.Context) (string, error)
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
	location, err := handleQueryLocation(ctx, controller)
	if err != nil {
		return entity.Weather{}, err
	}

	cookieId, _ := getIdFromCookie(ctx)
	// Not checking for error as new Id will be created in Db operations

	id, weatherData, err := controller.service.AddToRecent(location, cookieId)
	if err != nil {
		return entity.Weather{}, err
	}

	if cookieId != id {
		ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
	}

	return weatherData, nil
}

func (controller *controller) GetRecentsWeather(ctx *gin.Context) ([]entity.Weather, error) {
	cookieId, err := getIdFromCookie(ctx)
	if err != nil {
		return []entity.Weather{}, nil
	}

	recentsWeatherData, err := controller.service.GetRecentsWeather(cookieId)
	return recentsWeatherData, err
}

func (controller *controller) GetFavouritesWeather(ctx *gin.Context) ([]entity.Weather, error) {
	cookieId, err := getIdFromCookie(ctx)
	if err != nil {
		return []entity.Weather{}, nil
	}

	favouritesWeatherData, err := controller.service.GetFavouritesWeather(cookieId)
	return favouritesWeatherData, err
}

func (controller *controller) HandleFavourite(ctx *gin.Context) (string, error) {
	location, err := handleQueryLocation(ctx, controller)
	if err != nil {
		return "", err
	}

	cookieId, _ := getIdFromCookie(ctx)
	// Not checking for error as new Id will be created in Db operations

	id, response, err := controller.service.HandleFavourite(cookieId, location)
	if err != nil {
		return "", err
	}

	if cookieId != id {
		ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
	}

	return response, nil
}

func (controller *controller) ClearRecents(ctx *gin.Context) (string, error) {
	cookieId, err := getIdFromCookie(ctx)
	if err != nil {
		return "Recents are cleared", nil
	}

	err = controller.service.ClearRecents(cookieId)

	return "Recents are cleared", err
}

func (controller *controller) ClearFavourites(ctx *gin.Context) (string, error) {
	cookieId, err := getIdFromCookie(ctx)
	if err != nil {
		return "Favourites are cleared", nil
	}

	err = controller.service.ClearFavourites(cookieId)

	return "Favourites are cleared", err
}

// Get cookie, check it's validity and set if required
func getIdFromCookie(ctx *gin.Context) (string, error) {
	id, err := ctx.Cookie("id")

	// Id Cookie is not set
	if err != nil {
		return "", errors.New("Id cookie not set")
	}

	return id, nil
}

// Check for location in query params and send weather data
func handleQueryLocation(ctx *gin.Context, controller *controller) (string, error) {
	params := ctx.Request.URL.Query()

	if !params.Has("location") {
		return "", errors.New("No location")
	}

	location := params["location"][0]

	return location, nil
}
