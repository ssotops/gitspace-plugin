package gitspace_plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/charmbracelet/log"
)

// LoadPluginFunc is a function type for loading plugins
type LoadPluginFunc func(string) (GitspacePlugin, error)

// LoadPlugin is the default implementation of LoadPluginFunc
var LoadPlugin LoadPluginFunc = loadPluginImpl

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
