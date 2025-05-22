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
	
	// Initialize HTTP routes with CORS support
	http.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		
		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Process the actual request
		handleRecord(w, r)
	})

	// Start the server
	log.Println(fmt.Sprintf("Server starting on :%d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
