package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	configFile := "config.json"
	configs := []Config{}

	content, err := os.ReadFile(configFile)
	requestCount := 10

	for i := 0; i < requestCount; i++ {

		if err != nil {
			fmt.Println("Error reading file", err)
			return
		}

		err = json.Unmarshal(content, &configs)
		if err != nil {
			fmt.Println("Error unmarshalling JSON: ", err)
			return
		}

		for _, config := range configs {
			fmt.Printf("ID: %d, Method: %s, URL: %s\n ", config.ID, config.Payload.Method, config.Payload.URL)

			resp, err := http.Get(config.Payload.URL)
			if err != nil {
				fmt.Println("Error making HTTP request", err)
				continue
			}
			defer resp.Body.Close()

			fmt.Println("HTTP Response Status: ", resp.Status)
		}
	}

}
