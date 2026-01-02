package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// RuntimeRegistryPath returns the path to the global runtime registry config.
func RuntimeRegistryPath(homeDir string) string {
	return filepath.Join(homeDir, ".gastown", "runtimes.json")
}

// LoadRuntimeRegistryConfig loads the global runtime registry config.
func LoadRuntimeRegistryConfig(path string) (*RuntimeRegistryConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, path)
		}
		return nil, fmt.Errorf("reading runtime registry: %w", err)
	}

	var config RuntimeRegistryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing runtime registry: %w", err)
	}

	if config.Runtimes == nil {
		config.Runtimes = make(map[string]RuntimeAdapterConfig)
	}

	return &config, nil
}

// SaveRuntimeRegistryConfig saves the runtime registry config.
func SaveRuntimeRegistryConfig(path string, config *RuntimeRegistryConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding runtime registry: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing runtime registry: %w", err)
	}

	return nil
}
