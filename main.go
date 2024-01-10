package main

import (
	"encoding/json"
	"fmt"
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
	}

}
