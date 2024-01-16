package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
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

	for i := 0; i < requestCount; i++ {
		wg.Add(1)

		err = json.Unmarshal(content, &configs)
		if err != nil {
			fmt.Println("Error unmarshalling JSON: ", err)
			return
		}

		for _, config := range configs {
			fmt.Printf("ID: %d, Method: %s, URL: %s\n ", config.ID, config.Payload.Method, config.Payload.URL)

			if concurrent {
				go func(url string) {

					resp, err := http.Get(url)
					if err != nil {
						fmt.Println("Error making HTTP request", err)
						return
					}

					defer resp.Body.Close()

					fmt.Println("HTTP Response Status: ", resp.Status)
				}(config.Payload.URL)
			} else {

				resp, err := http.Get(config.Payload.URL)
				if err != nil {
					fmt.Println("Error making HTTP request", err)
					return
				}

				defer resp.Body.Close()

				fmt.Println("HTTP Response Status: ", resp.Status)
			}

		}
	}

	if concurrent {
		wg.Wait()
	}

}
