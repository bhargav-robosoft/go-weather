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

	id, err := ctx.Cookie("id")

	if err != nil {
		id = db.CreateIdForUser()
		ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
	}

	err = db.AddRecentSearchForUser(id, data.Name)
	if err != nil {
		if err.Error() == "Id error" {
			id := db.CreateIdForUser()
			ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
			err := db.AddRecentSearchForUser(id, data.Name)
			if err != nil {
				return entity.Weather{}, err
			}
		} else {
			return entity.Weather{}, err
		}

	}

	isFav, err := db.IsFavourite(id, data.Name)

	if err != nil {
		return data, nil
	}

	data.IsFavourite = isFav

	return data, nil
}

func (controller *controller) GetRecentsWeather(ctx *gin.Context) ([]entity.Weather, error) {
	id, err := ctx.Cookie("id")
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
	id, err := ctx.Cookie("id")
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
	params := ctx.Request.URL.Query()

	if !params.Has("location") {
		return "", errors.New("No location")
	}

	location := params["location"][0]

	if len(location) == 0 {
		return "", errors.New("Empty location")
	}

	data, err := controller.service.GetWeather(location)
	if err != nil {
		return "", err
	}

	cookie, err := ctx.Cookie("id")

	if err != nil {
		id := db.CreateIdForUser()
		ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
		cookie = id
	}

	response, err := db.HandleFavouriteForUser(cookie, data.Name)
	if err != nil {
		if err.Error() == "Invalid Id" || err.Error() == "Id not found" {
			id := db.CreateIdForUser()
			ctx.SetCookie("id", id, 3600, "/", ctx.Request.Host, true, true)
			response, err := db.HandleFavouriteForUser(id, data.Name)
			if err != nil {
				return "", err
			} else {
				return response, nil
			}
		} else {
			return "", err
		}
	}
	return response, nil
}
