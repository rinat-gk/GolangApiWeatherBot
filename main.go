package main

const OpenWeatherMapAPIKey = "0e264e7da407afac79c77c3f2d1feeeb"

type WeatherMapApiResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}
