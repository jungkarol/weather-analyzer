package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
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

	citiesChan := make(chan models.City, len(cities))
	for _, city := range cities {
		citiesChan <- city
	}
	close(citiesChan)

	weatherDataChan := make(chan models.WeatherData, len(cities))

	var producersWg sync.WaitGroup
	for i := 0; i < config.Producers; i++ {
		producersWg.Add(1)
		go func(producerId int) {
			defer producersWg.Done()
			producer(producerId, citiesChan, weatherDataChan)
		}(i)
	}

	go func() {
		producersWg.Wait()
		close(weatherDataChan)
	}()

	startTime := time.Now()

	var consumersWg sync.WaitGroup
	finalResults := models.Results{}
	var resultsMutex sync.Mutex

	for i := 0; i < config.Consumers; i++ {
		consumersWg.Add(1)
		go func(consumerId int) {
			defer consumersWg.Done()
			consumer(consumerId, weatherDataChan, &finalResults, &resultsMutex)
		}(i)
	}

	consumersWg.Wait()

	fmt.Printf("Czas przetwarzania przez konsumentów: %v\n", time.Since(startTime))
	saveResultsToFile(finalResults)
}
