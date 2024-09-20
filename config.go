package gitspace_plugin

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
)

func ParsePluginConfig(pluginDir string) (PluginConfig, error) {
	var config PluginConfig
	configPath := filepath.Join(pluginDir, "gitspace-plugin.toml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.Metadata.Name == "" {
		return config, fmt.Errorf("plugin name is required")
	}
	if config.Metadata.Version == "" {
		config.Metadata.Version = "0.1.0" // Set a default version if not specified
	}

	return config, nil
}
