package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"weather-analyzer/models"
)

func consumer(ID int, weatherDataChan <-chan models.WeatherData, finalResults *models.Results, resultsMutex *sync.Mutex) {
	for weatherData := range weatherDataChan {
		startTime := time.Now()
		avgTemp := calculateAverageTemperature(&weatherData)
		foggyDays := countFoggyDays(&weatherData)
		sunnyDays := countSunnyDays(&weatherData)

		computedData := models.WeatherComputedData{
			City:      weatherData.City,
			AvgTemp:   avgTemp,
			FoggyDays: foggyDays,
			SunnyDays: sunnyDays,
		}

		resultsMutex.Lock()

		if computedData.AvgTemp > finalResults.MaxTemperature {
			finalResults.MaxTemperature = computedData.AvgTemp
			finalResults.HottestCity = computedData.City.Name
		}

		if computedData.FoggyDays > finalResults.MaxFoggyDays {
			finalResults.MaxFoggyDays = computedData.FoggyDays
			finalResults.FoggiestCity = computedData.City.Name
		}

		if computedData.SunnyDays > finalResults.MaxSunnyDays {
			finalResults.MaxSunnyDays = computedData.SunnyDays
			finalResults.ClearestCity = computedData.City.Name
		}

		resultsMutex.Unlock()

		elapsedTime := float64(time.Since(startTime).Nanoseconds()) / 1_000_000.0
		fmt.Printf("Consumer [%d] Przetwarzanie %s zakończone w %.5f ms\n", ID, weatherData.City.Name, elapsedTime)
	}
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
