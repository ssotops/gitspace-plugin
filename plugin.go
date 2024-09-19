package gitspace_plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/charmbracelet/log"
)

// LoadPluginFunc is a function type for loading plugins
type LoadPluginFunc func(string) (GitspacePlugin, error)

// LoadPlugin is the default implementation of LoadPluginFunc
var LoadPlugin LoadPluginFunc = loadPluginImpl

var SharedDependencies map[string]string

func SetSharedDependencies(deps map[string]string) {
	SharedDependencies = deps
}

func GetSharedDependencies() map[string]string {
	return SharedDependencies
}

// loadPluginImpl is the actual implementation of plugin loading
func loadPluginImpl(pluginPath string) (GitspacePlugin, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin does not have a Plugin symbol: %w", err)
	}

	plugin, ok := symPlugin.(GitspacePlugin)
	if !ok {
		return nil, fmt.Errorf("plugin does not implement GitspacePlugin interface")
	}

	return plugin, nil
}

// LoadPluginWithConfig loads a plugin and sets its configuration
func LoadPluginWithConfig(pluginPath string) (GitspacePlugin, error) {
	plugin, err := LoadPlugin(pluginPath)
	if err != nil {
		return nil, err
	}

	// Check and update dependencies
	if err := updatePluginDependencies(plugin, pluginPath); err != nil {
		return nil, fmt.Errorf("failed to update plugin dependencies: %w", err)
	}

	pluginDir := filepath.Dir(pluginPath)
	config, err := ParsePluginConfig(pluginDir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plugin config: %w", err)
	}

	plugin.SetConfig(config)
	return plugin, nil
}

// RunPlugin runs the given plugin
func RunPlugin(plugin GitspacePlugin, logger *log.Logger) error {
	logger.Info("Running plugin", "name", plugin.Name(), "version", plugin.Version())
	return plugin.Run(logger)
}

// RunStandalonePlugin runs the plugin in standalone mode
func RunStandalonePlugin(plugin GitspacePlugin, args []string) error {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.InfoLevel)
	logger.Info("Running plugin in standalone mode", "name", plugin.Name(), "version", plugin.Version())
	return plugin.Standalone(args)
}

func updatePluginDependencies(plugin GitspacePlugin, pluginPath string) error {
	pluginDeps := plugin.GetDependencies()
	sharedDeps := GetSharedDependencies()

	needsUpdate := false
	for dep, version := range sharedDeps {
		if pluginDeps[dep] != version {
			needsUpdate = true
			break
		}
	}

	if needsUpdate {
		// Update go.mod
		cmd := exec.Command("go", "mod", "edit")
		for dep, version := range sharedDeps {
			cmd.Args = append(cmd.Args, "-require", fmt.Sprintf("%s@%s", dep, version))
		}
		cmd.Dir = filepath.Dir(pluginPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to update go.mod: %w\nOutput: %s", err, output)
		}

		// Rebuild plugin
		if err := rebuildPlugin(pluginPath); err != nil {
			return fmt.Errorf("failed to rebuild plugin: %w", err)
		}

		// Reload plugin
		newPlugin, err := LoadPlugin(pluginPath)
		if err != nil {
			return fmt.Errorf("failed to reload updated plugin: %w", err)
		}
		// Instead of assigning pointers, update the interface
		if updatablePlugin, ok := plugin.(interface{ Update(GitspacePlugin) }); ok {
			updatablePlugin.Update(newPlugin)
		} else {
			return fmt.Errorf("plugin does not support updating")
		}
	}

	return nil
}

func rebuildPlugin(pluginPath string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginPath)
	cmd.Dir = filepath.Dir(pluginPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to rebuild plugin: %w\nOutput: %s", err, output)
	}
	return nil
}
