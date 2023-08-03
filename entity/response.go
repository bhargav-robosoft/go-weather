package entity

type WeatherSuccessResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    Weather `json:"data"`
}

type WeatherFailureResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Example string `json:"example"`
}

type MultiWeatherResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Weather `json:"data"`
}
