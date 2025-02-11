package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"weather-analyzer/config"
	"weather-analyzer/models"
)

func producer(cities []models.City, weatherDataChan chan<- models.WeatherData) {
	for _, city := range cities {
		startTime := time.Now()
		weatherData, err := getWeatherData(city)
		if err != nil {
			fmt.Printf("Błąd pobierania dla %s: %v\n", city.Name, err)
			continue
		}
		elapsedTime := time.Since(startTime)
		fmt.Printf("Pobranie zakończone dla %s w %s\n", city.Name, elapsedTime)
		weatherData.City = city
		weatherDataChan <- *weatherData
	}
}

func getWeatherData(city models.City) (*models.WeatherData, error) {
	startDate := time.Now().AddDate(0, config.MonthsBack*-1, 0).Format("2006-01-02") // 6 miesięcy wstecz
	endDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")                     // Wczorajsza data

	baseUrl := config.WeatherAPIBaseURL
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&start_date=%s&end_date=%s&daily=temperature_2m_max,temperature_2m_min,weathercode&timezone=Europe/Warsaw", baseUrl, city.Latitude, city.Longitude, startDate, endDate)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("błąd pobierania API: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var weatherResponse models.WeatherData
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		fmt.Printf("Response body: %s", body)
		return nil, fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	return &weatherResponse, nil
}
