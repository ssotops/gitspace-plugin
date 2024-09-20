package gsplug

import (
	"fmt"
  "io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
  "strings"
)

// BuildPlugin builds the plugin in the specified directory
func BuildPlugin(pluginDir string) error {
	// Ensure the plugin directory exists
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin directory does not exist: %s", pluginDir)
	}

	// Read the plugin manifest
	manifest, err := ReadManifest(filepath.Join(pluginDir, "gitspace-plugin.toml"))
	if err != nil {
		return fmt.Errorf("failed to read plugin manifest: %w", err)
	}

	// Check compatibility
	compatible, err := CheckCompatibility(manifest.Metadata.Version)
	if err != nil {
		return fmt.Errorf("failed to check compatibility: %w", err)
	}
	if !compatible {
		return fmt.Errorf("plugin version %s is not compatible with the current Gitspace version", manifest.Metadata.Version)
	}

	// Update dependencies
	if err := UpdatePluginDependencies(pluginDir); err != nil {
		return fmt.Errorf("failed to update plugin dependencies: %w", err)
	}

	// Get the plugin name from the directory name
	pluginName := filepath.Base(pluginDir)

	canonicalDeps, err := GetCanonicalDeps()
	if err != nil {
		return fmt.Errorf("failed to get canonical dependencies: %w", err)
	}

	// Update go.mod with canonical versions
	if err := updateGoMod(pluginDir, canonicalDeps); err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// Build the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", filepath.Join(pluginDir, "dist", pluginName+".so"))
	cmd.Dir = pluginDir
	cmd.Env = append(os.Environ(), "GOPROXY=direct")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

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

func updateGoMod(pluginDir string, canonicalDeps CanonicalDeps) error {
    goModPath := filepath.Join(pluginDir, "go.mod")
    content, err := ioutil.ReadFile(goModPath)
    if err != nil {
        return err
    }

    lines := strings.Split(string(content), "\n")
    var newLines []string
    for _, line := range lines {
        if strings.HasPrefix(line, "require ") {
            parts := strings.Fields(line)
            if len(parts) >= 3 {
                module := parts[1]
                if version, ok := canonicalDeps.Versions[module]; ok {
                    newLines = append(newLines, fmt.Sprintf("require %s %s", module, version))
                    continue
                }
            }
        }
        newLines = append(newLines, line)
    }

    return ioutil.WriteFile(goModPath, []byte(strings.Join(newLines, "\n")), 0644)
}
