package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// FridgeItem represents a single item to be stored in the fridge
type FridgeItem struct {
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	Unit       string `json:"unit"`
	ExpireDate string `json:"expireDate"`
}

// FridgeResponse represents the structured response
type FridgeResponse struct {
	Items []FridgeItem `json:"items"`
}

// RecordRequest represents the incoming request
type RecordRequest struct {
	Description string `json:"description"`
	CurrentTime string `json:"currentTime,omitempty"` // Optional, format: YYYY-MM-DD HH:MM:SS
}

// handleRecord processes the natural language description and returns structured data
func handleRecord(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the request
	var req RecordRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error parsing JSON request", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Description == "" {
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	// Process the description with Gemini API
	response, err := processWithGemini(req.Description, req.CurrentTime)
	if err != nil {
		log.Printf("Error processing with Gemini: %v", err)
		http.Error(w, "Error processing with AI", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Return the structured response
	json.NewEncoder(w).Encode(response)
}
