package controller

import (
	"errors"
	"weather-app/db"
	"weather-app/entity"
	"weather-app/service"

	"github.com/gin-gonic/gin"
)

type WeatherController interface {
	GetWeather(ctx *gin.Context) (entity.Weather, error)
	GetRecentsWeather(ctx *gin.Context) ([]entity.Weather, error)
	GetFavouritesWeather(ctx *gin.Context) ([]entity.Weather, error)
	HandleFavourite(ctx *gin.Context) (string, error)
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
	weatherData, err := handleQueryLocation(ctx, controller)
	if err != nil {
		return entity.Weather{}, err
	}

	id, err := handleCookie(ctx, true)
	if err != nil {
		return entity.Weather{}, err
	}

	err = db.AddRecentSearchForUser(id, weatherData.Name)
	if err != nil {
		return entity.Weather{}, err
	}

	isFav, err := db.IsFavourite(id, weatherData.Name)

	if err != nil {
		return weatherData, nil
	}

	weatherData.IsFavourite = isFav

	return weatherData, nil
}

func (controller *controller) GetRecentsWeather(ctx *gin.Context) ([]entity.Weather, error) {
	id, err := handleCookie(ctx, false)
	if err != nil {
		return []entity.Weather{}, nil
	} else {
		recents, err := db.GetRecentsForUser(id)

		if err != nil {
			return []entity.Weather{}, err
		} else {
			recentWeatherData, err := controller.service.GetRecents(recents)
			if err != nil {
				return []entity.Weather{}, err
			} else {
				for i, d := range recentWeatherData {
					isFav, err := db.IsFavourite(id, d.Name)
					if err != nil {
						continue
					}
					recentWeatherData[i].IsFavourite = isFav
				}
				return recentWeatherData, nil
			}
		}
	}
}

func (controller *controller) GetFavouritesWeather(ctx *gin.Context) ([]entity.Weather, error) {
	id, err := handleCookie(ctx, false)
	if err != nil {
		return []entity.Weather{}, nil
	} else {
		favourites, err := db.GetFavouritesForUser(id)

		if err != nil {
			return []entity.Weather{}, err
		} else {
			favouriteWeatherData, err := controller.service.GetFavourites(favourites)
			if err != nil {
				return []entity.Weather{}, err
			} else {
				for i := range favouriteWeatherData {
					favouriteWeatherData[i].IsFavourite = true
				}
				return favouriteWeatherData, nil
			}
		}
	}
}

func (controller *controller) HandleFavourite(ctx *gin.Context) (string, error) {
	weatherData, err := handleQueryLocation(ctx, controller)
	if err != nil {
		return "", err
	}

	id, err := handleCookie(ctx, true)
	if err != nil {
		return "", err
	}

	response, err := db.HandleFavouriteForUser(id, weatherData.Name)
	if err != nil {
		return "", err
	}

	return response, nil
}

// Get cookie, check it's validity and set if required
func handleCookie(ctx *gin.Context, setId bool) (string, error) {
	id, err := ctx.Cookie("id")

	// Id Cookie is not set
	if err != nil {
		if setId {
			id = db.CreateIdForUser()
			ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
			return id, nil
		} else {
			return "", errors.New("Id cookie not set")
		}
	}

	_, err = db.CheckUserId(id)

	// Id is invalid
	if err != nil {
		if setId {
			id = db.CreateIdForUser()
			ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
		} else {
			return "", errors.New("Invalid Id")
		}
	}

	return id, nil
}

// Check for location in query params and send weather data
func handleQueryLocation(ctx *gin.Context, controller *controller) (entity.Weather, error) {
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
