package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
	gitspace_plugin "github.com/ssotops/gitspace-plugin"
)

type HelloWorldPlugin struct{}

var Plugin HelloWorldPlugin

func (p HelloWorldPlugin) Name() string {
	return "hello-world"
}

func (p HelloWorldPlugin) Version() string {
	return "1.0.0"
}

func (p HelloWorldPlugin) Description() string {
	return "A simple Hello World plugin for Gitspace"
}

func (p HelloWorldPlugin) Run(logger *log.Logger) error {
	logger.Info("Hello from the Hello World plugin!")
	return nil
}

func (p HelloWorldPlugin) GetMenuOption() *gitspace_plugin.Option {
	return &gitspace_plugin.Option{
		Key:   "hello-world",
		Value: "Hello World",
	}
}

func (p HelloWorldPlugin) Standalone(args []string) error {
	fmt.Println("Hello from the standalone Hello World plugin!")
	return nil
}

func (p HelloWorldPlugin) SetConfig(config gitspace_plugin.PluginConfig) {
	// This plugin doesn't use any configuration, but we need to implement this method
}

func (p HelloWorldPlugin) GetDependencies() map[string]string {
	return gitspace_plugin.GetSharedDependencies()
}

func (p HelloWorldPlugin) Update(newPlugin gitspace_plugin.GitspacePlugin) {
	// This method is called when the plugin is updated
	// For this simple plugin, we don't need to do anything
}

func main() {
	plugin := HelloWorldPlugin{}
	if err := plugin.Standalone(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// The following functions are not part of the plugin interface
// They are here for demonstration purposes and to fix the build errors

func updatePluginDependencies(plugin gitspace_plugin.GitspacePlugin, pluginPath string) error {
	pluginDeps := plugin.GetDependencies()
	sharedDeps := gitspace_plugin.GetSharedDependencies()

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
		newPlugin, err := gitspace_plugin.LoadPlugin(pluginPath)
		if err != nil {
			return fmt.Errorf("failed to reload updated plugin: %w", err)
		}
		plugin.Update(newPlugin)
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
