package config

import (
	"encoding/json"
	"os"
)

// Add this to allow easy mocking
var (
	OpenFile   = os.Open
	NewDecoder = json.NewDecoder
)

type Config struct {
	SomeConfig string `json:"config_data"`
}

func LoadConfig() *Config {
	file, err := OpenFile("config/config.json")
	if err != nil {
		return &Config{}
	}

	defer file.Close()

	var config Config
	decoder := NewDecoder(file)

	// Add additional error handling
	err = decoder.Decode(&config)
	if err != nil {
		// Log the error if needed
		return &Config{}
	}

	return &config
}
