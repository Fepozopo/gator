package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config represents the JSON structure of the configuration file.
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// getConfigFilePath constructs the full path to the configuration file.
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error obtaining home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

// write writes the Config struct to the JSON configuration file.
func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Error obtaining config filepath: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(cfg)
}

// Read reads the configuration from the JSON file and returns a Config struct.
func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error obtaining config filepath: %w", err)
	}

	file, err := os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		// Return an empty Config if the file does not exist.
		return Config{}, nil
	} else if err != nil {
		return Config{}, fmt.Errorf("Error opening config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("Error decoding config file: %w", err)
	}
	return cfg, nil
}

// SetUser updates the current_user_name field and writes the updated Config back to the file.
func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}
