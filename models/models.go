package models

type City struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type WeatherComputedData struct {
	City      City    `json:"city"`
	AvgTemp   float64 `json:"avg_temp"`
	FoggyDays int     `json:"foggy_days"`
	SunnyDays int     `json:"sunny_days"`
}

type Results struct {
	HottestCity  string `json:"hottest_city"`
	FoggiestCity string `json:"foggiest_city"`
	ClearestCity string `json:"clearest_city"`
}

type WeatherData struct {
	City  City `json:"city"`
	Daily struct {
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		WeatherCode    []int     `json:"weathercode"`
	} `json:"daily"`
}
