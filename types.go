package gitspace_plugin

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

// GitspaceCatalog struct needs to be defined
type GitspaceCatalog struct {
	Plugins map[string]CatalogPlugin
}

type CatalogPlugin struct {
	Path string
	// Add other necessary fields
}

type GitspacePlugin interface {
	Name() string
	Version() string
	Description() string
	Run(logger *log.Logger) error
	GetMenuOption() *huh.Option[string]
	Standalone(args []string) error
	SetConfig(config PluginConfig)
	GetDependencies() map[string]string
	Update(GitspacePlugin)
}

// TODO: Refactor and simplify types below

// PluginMetadata contains additional information about the plugin
type PluginMetadata struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Description string   `toml:"description"`
	Author      string   `toml:"author"`
	Tags        []string `toml:"tags"`
}

type Option = huh.Option[string]

// PluginConfig contains the configuration for the plugin
type PluginConfig struct {
	Metadata PluginMetadata `toml:"metadata"`
	Menu     struct {
		Title string `toml:"title"`
		Key   string `toml:"key"`
	} `toml:"menu"`
}

type PluginManifest struct {
	Metadata struct {
		Name        string `toml:"name"`
		Version     string `toml:"version"`
		Description string `toml:"description"`
		Author      string `toml:"author"`
	} `toml:"metadata"`
	Menu struct {
		Title string `toml:"title"`
		Key   string `toml:"key"`
	} `toml:"menu"`
	Sources []struct {
		Path       string `toml:"path"`
		EntryPoint string `toml:"entry_point,omitempty"`
	} `toml:"sources"`
}
