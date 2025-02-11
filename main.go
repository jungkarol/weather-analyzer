package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"weather-analyzer/config"
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
	resultsChan := make(chan models.Results)

	go func() {
		producer(cities, weatherDataChan)
		close(weatherDataChan)
	}()

	var consumersWg sync.WaitGroup
	for i := 0; i < config.Consumers; i++ {
		consumersWg.Add(1)
		go func(consumerId int) {
			defer consumersWg.Done()
			consumer(consumerId, weatherDataChan, resultsChan)
		}(i)
	}

	go func() {
		consumersWg.Wait()
		close(resultsChan)
	}()

	var results models.Results
	for r := range resultsChan {
		results = r
	}

	saveResultsToFile(results)

	fmt.Println("Proces zakończony!")
}
