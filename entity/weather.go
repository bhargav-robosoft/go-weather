package entity

type Weather struct {
	Name            string  `json:"name"`
	CountryName     string  `json:"country-name"`
	Temperature     float64 `json:"temperature"`
	Description     string  `json:"description"`
	WeatherIconLink string  `json:"icon"`
	MinTemperature  float64 `json:"minimum-temperature"`
	MaxTemperature  float64 `json:"maximum-temperature"`
	Clouds          int     `json:"clouds"`
	Humidity        int     `json:"humidity"`
	WindSpeed       float64 `json:"wind-speed"`
	Visibility      int     `json:"visibility"`
}

type WeatherError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Example string `json:"example"`
}
