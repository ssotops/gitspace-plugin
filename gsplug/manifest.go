package gsplug

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

func ReadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest PluginManifest
	err = toml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func WriteManifest(manifest *PluginManifest, path string) error {
	data, err := toml.Marshal(manifest)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
