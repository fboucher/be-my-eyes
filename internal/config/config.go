package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	APIKey string `json:"api_key"`
}

// configDir returns the configuration directory path
// This is ~/.config/be-my-eyes
func configDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "be-my-eyes"), nil
}

// configFilePath returns the full path to the config file
func configFilePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads the configuration from the config file
// Returns an empty config if the file doesn't exist
func Load() (*Config, error) {
	path, err := configFilePath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return empty config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path, err := configFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// EnsureAPIKey checks if an API key is configured and prompts for one if not
func EnsureAPIKey() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	if cfg.APIKey != "" {
		return cfg.APIKey, nil
	}

	// If no API key, check environment variable
	if apiKey := os.Getenv("REKA_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
		if err := cfg.Save(); err != nil {
			return "", fmt.Errorf("failed to save API key from environment: %w", err)
		}
		return apiKey, nil
	}

	return "", fmt.Errorf("no API key found. Please set REKA_API_KEY environment variable or add it to the config file")
}
