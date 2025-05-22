package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	GeminiAPIKey string `yaml:"gemini_api_key"`
	GeminiModel  string `yaml:"gemini_model"`
}

// Global configuration instance
var (
	appConfig     *Config
	configOnce    sync.Once
	configLoadErr error
)

// GetConfig returns the cached configuration or loads it if not yet loaded
func GetConfig() (*Config, error) {
	configOnce.Do(func() {
		appConfig, configLoadErr = loadConfigFromFile()
		if configLoadErr != nil {
			log.Printf("Warning: Failed to load config: %v, using defaults and environment variables", configLoadErr)
			// Even on error, initialize with defaults
			appConfig = &Config{
				GeminiModel:  "gemini-pro",
				GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
			}
			configLoadErr = nil // Reset error since we're using defaults
		}
		log.Printf("Configuration loaded, using model: %s", appConfig.GeminiModel)
	})
	return appConfig, configLoadErr
}

// loadConfigFromFile loads configuration from config.yaml file
func loadConfigFromFile() (*Config, error) {
	// Default configuration
	config := &Config{
		GeminiModel: "gemini-pro", // Default model
	}

	// Try to find config file
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If config file doesn't exist in current directory, check for it in the executable directory
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			configPath = filepath.Join(execDir, "config.yaml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// If still not found, just return default config with environment variable fallback
				config.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
				return config, nil
			}
		} else {
			// If we can't determine executable path, just return default config
			config.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
			return config, nil
		}
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// If API key is not set in config, try environment variable as fallback
	if config.GeminiAPIKey == "" {
		config.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	}

	// If model is not set, use default
	if config.GeminiModel == "" {
		config.GeminiModel = "gemini-pro"
	}

	return config, nil
}
