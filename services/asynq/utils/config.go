package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ClientConfig represents configuration for a gRPC client
type ClientConfig struct {
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Timeout int               `json:"timeout"`
	Options map[string]string `json:"options"`
}

// ServiceConfig represents the overall service configuration
type ServiceConfig struct {
	Clients []ClientConfig `json:"clients"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(configPath string) (*ServiceConfig, error) {
	// Default config path
	if configPath == "" {
		configPath = os.Getenv("ASYNQ_CONFIG_PATH")
		if configPath == "" {
			configPath = "./config.json"
		}
	}

	// Read config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		// Return default config if file doesn't exist
		return getDefaultConfig(), nil
	}

	// Parse config
	var config ServiceConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// getDefaultConfig returns default configuration
func getDefaultConfig() *ServiceConfig {
	return &ServiceConfig{
		Clients: []ClientConfig{
			{
				Name:    "auth-service",
				Address: "localhost:50052",
				Timeout: 30,
			},
			{
				Name:    "user-service",
				Address: "localhost:50053",
				Timeout: 30,
			},
		},
	}
}