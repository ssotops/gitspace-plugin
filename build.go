package gsplug

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildPlugin builds the plugin in the specified directory
func BuildPlugin(pluginDir string) error {
	// Ensure the plugin directory exists
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin directory does not exist: %s", pluginDir)
	}

	// Update dependencies
	if err := UpdatePluginDependencies(pluginDir); err != nil {
		return fmt.Errorf("failed to update plugin dependencies: %w", err)
	}

	// Get the plugin name from the directory name
	pluginName := filepath.Base(pluginDir)

	// Build the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", filepath.Join(pluginDir, pluginName+".so"))
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	fmt.Printf("Successfully built plugin: %s\n", pluginName)
	return nil
}

// BuildAllPlugins builds all plugins in the Gitspace plugins directory
func BuildAllPlugins() error {
	pluginsDir := filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", "plugins")

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pluginDir := filepath.Join(pluginsDir, entry.Name())
			if err := BuildPlugin(pluginDir); err != nil {
				fmt.Printf("Failed to build plugin %s: %v\n", entry.Name(), err)
			}
		}
	}

	return nil
}
