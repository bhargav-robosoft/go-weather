package entity

type Weather struct {
	Name            string `json:"name"`
	CountryName     string `json:"country-name"`
	Temperature     string `json:"temperature"`
	Description     string `json:"description"`
	WeatherIconLink string `json:"icon"`
	MinTemperature  string `json:"minimum-temperature"`
	MaxTemperature  string `json:"maximum-temperature"`
	Clouds          string `json:"clouds"`
	Humidity        string `json:"humidity"`
	WindSpeed       string `json:"wind-speed"`
	Visibility      string `json:"visibility"`
}

type WeatherError struct {
	Status  int    `json:"status"`
	Message string `json:"location"`
	Example string `json:"eample"`
}
