package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8602
	
	// Initialize configuration
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	log.Printf("Configuration initialized, using Gemini model: %s", config.GeminiModel)
	
	// Initialize HTTP routes
	http.HandleFunc("/record", handleRecord)

	// Start the server
	log.Println(fmt.Sprintf("Server starting on :%d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
