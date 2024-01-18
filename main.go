package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type RequestPayload struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Config struct {
	ID      int            `json:"ID"`
	Payload RequestPayload `json:"Payload"`
}

type Statistics struct {
	TotalRequests      int
	FailedRequests     int
	SuccessfulRequests int
	RequestsPerSecond  float64
}

var statistics = Statistics{}

func main() {

	var configFile string
	var concurrent bool
	var requestCount int

	flag.StringVar(&configFile, "config", "config.json", "Path to the configuration file")
	flag.BoolVar(&concurrent, "concurrent", true, "Whether to make requests concurrently")
	flag.IntVar(&requestCount, "requests", 10, "Number of requests to make")

	flag.Parse()

	configs := []Config{}

	content, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < requestCount; i++ {
		wg.Add(1)

		err = json.Unmarshal(content, &configs)
		if err != nil {
			fmt.Println("Error unmarshalling JSON: ", err)
			return
		}

		for _, config := range configs {
			if concurrent {
				go func(url string) {
					defer wg.Done()
					makeRequest(url)
				}(config.Payload.URL)
			} else {
				makeRequest(config.Payload.URL)
			}
		}
	}

	if concurrent {
		wg.Wait()
	}

	elapsedTime := time.Since(startTime).Seconds()
	printStatistics(elapsedTime)

}

func makeRequest(url string) {
	statistics.TotalRequests++

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request", err)
		statistics.FailedRequests++
		return
	}

	defer resp.Body.Close()

	statisticsUpdateTimings(resp.StatusCode)
}

func printStatistics(elapsedTime float64) {

	statistics.RequestsPerSecond = float64(statistics.TotalRequests) / elapsedTime

	fmt.Println("Results ...")
	fmt.Printf("  Total Requests .......................: %d\n", statistics.TotalRequests)
	fmt.Printf("  Total Successful Requests (2XX).......................: %d\n", statistics.SuccessfulRequests)
	fmt.Printf("  Failed Requests (5XX).......................: %d\n", statistics.FailedRequests)
	fmt.Printf("  Requests Per Second.......................: %.2f\n", statistics.RequestsPerSecond)
}

func statisticsUpdateTimings(statusCode int) {
	if statusCode >= 200 && statusCode < 300 {
		statistics.SuccessfulRequests++
	}
}
