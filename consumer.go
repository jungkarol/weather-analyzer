package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"weather-analyzer/models"
)

func consumer(weatherDataChan <-chan models.WeatherData, resultsChan chan<- models.Results, wg *sync.WaitGroup) {
	defer wg.Done()

	var hottestCity models.WeatherComputedData
	var foggiestCity models.WeatherComputedData
	var clearestCity models.WeatherComputedData

	for weatherData := range weatherDataChan {
		startTime := time.Now()
		avgTemp := calculateAverageTemperature(&weatherData)
		foggyDays := countFoggyDays(&weatherData)
		sunnyDays := countSunnyDays(&weatherData)
		elapsedTime := float64(time.Since(startTime).Nanoseconds()) / 1_000_000.0
		fmt.Printf("Przetwarzanie %s zakończone w %.5f ms\n", weatherData.City.Name, elapsedTime)
		fmt.Printf("For %s avgTemp: %f foggyDays: %d sunnyDays: %d\n", weatherData.City.Name, avgTemp, foggyDays, sunnyDays)

		computedData := models.WeatherComputedData{
			City:      weatherData.City,
			AvgTemp:   avgTemp,
			FoggyDays: foggyDays,
			SunnyDays: sunnyDays,
		}

		if computedData.AvgTemp > hottestCity.AvgTemp {
			hottestCity = computedData
		}

		if computedData.FoggyDays > foggiestCity.FoggyDays {
			foggiestCity = computedData
		}

		if computedData.SunnyDays > clearestCity.SunnyDays {
			clearestCity = computedData
		}
	}

	results := models.Results{
		HottestCity:  hottestCity.City.Name,
		FoggiestCity: foggiestCity.City.Name,
		ClearestCity: clearestCity.City.Name,
	}

	resultsChan <- results
}

func saveResultsToFile(results models.Results) {
	filename := "results.json"
	deleteFileIfExists(filename)

	file, _ := os.Create(filename)
	defer file.Close()

	json.NewEncoder(file).Encode(results)
	fmt.Println("Wyniki zapisane do results.json")
}

func deleteFileIfExists(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
	}
}

func calculateAverageTemperature(weatherData *models.WeatherData) float64 {
	total := 0.0
	count := 0

	for i := range weatherData.Daily.TemperatureMax {
		temp := (weatherData.Daily.TemperatureMax[i] + weatherData.Daily.TemperatureMin[i]) / 2
		total += temp
		count++
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func countFoggyDays(weatherData *models.WeatherData) int {
	count := 0
	for _, code := range weatherData.Daily.WeatherCode {
		if code == 45 {
			count++
		}
	}
	return count
}

func countSunnyDays(weatherData *models.WeatherData) int {
	count := 0
	for _, code := range weatherData.Daily.WeatherCode {
		if code == 0 {
			count++
		}
	}
	return count
}
