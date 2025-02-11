package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"weather-analyzer/models"
)

func main() {
	file, _ := os.ReadFile("resources/pl172.json")
	var cities []models.City
	err := json.Unmarshal(file, &cities)
	if err != nil {
		return
	}

	weatherDataChan := make(chan models.WeatherData, len(cities))
	resultsChan := make(chan models.Results, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		producer(cities, weatherDataChan, &wg)
		close(weatherDataChan)
	}()

	wg.Add(1)
	go func() {
		consumer(weatherDataChan, resultsChan, &wg)
		close(resultsChan)
	}()

	wg.Wait()

	results := <-resultsChan

	saveResultsToFile(results)

	fmt.Println("Proces zakończony!")
}
