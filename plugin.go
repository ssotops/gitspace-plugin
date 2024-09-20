package gitspace_plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
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

func GetCurrentDependencies() (map[string]string, error) {
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var modules []struct {
		Path    string
		Version string
	}

	decoder := json.NewDecoder(bytes.NewReader(output))
	for decoder.More() {
		var module struct {
			Path    string
			Version string
		}
		if err := decoder.Decode(&module); err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	dependencies := make(map[string]string)
	for _, module := range modules {
		dependencies[module.Path] = module.Version
	}

	return dependencies, nil
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
		// If there's an error loading the plugin, try rebuilding it
		if rebuildErr := rebuildAndLoadPlugin(pluginPath); rebuildErr != nil {
			return nil, fmt.Errorf("failed to load plugin and rebuild failed: %w", rebuildErr)
		}
		// Try loading the plugin again after rebuilding
		plugin, err = LoadPlugin(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load plugin after rebuild: %w", err)
		}
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

func rebuildAndLoadPlugin(pluginPath string) error {
	pluginDir := filepath.Dir(pluginPath)

	// Update go.mod to use the same versions as Gitspace
	if err := updateGoMod(pluginDir); err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// Rebuild the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginPath)
	cmd.Dir = pluginDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rebuild plugin: %w\nOutput: %s", err, output)
	}

	return nil
}

func updateGoMod(pluginDir string) error {
	// Get Gitspace's dependencies
	gitspaceDeps, err := GetCurrentDependencies()
	if err != nil {
		return fmt.Errorf("failed to get Gitspace dependencies: %w", err)
	}

	// Update go.mod
	for dep, version := range gitspaceDeps {
		cmd := exec.Command("go", "get", fmt.Sprintf("%s@%s", dep, version))
		cmd.Dir = pluginDir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to update dependency %s: %w\nOutput: %s", dep, err, output)
		}
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = pluginDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to tidy go.mod: %w\nOutput: %s", err, output)
	}

	return nil
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
			if version != "" {
				cmd.Args = append(cmd.Args, "-require", fmt.Sprintf("%s@%s", dep, version))
			}
		}
		cmd.Dir = filepath.Dir(pluginPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to update go.mod: %w\nOutput: %s", err, output)
		}

		// Run go mod tidy
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = filepath.Dir(pluginPath)
		if output, err := tidyCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to tidy go.mod: %w\nOutput: %s", err, output)
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

		// Update the plugin interface
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

// UpdatePluginDependencies is the public function to update a plugin's dependencies
func UpdatePluginDependencies(plugin GitspacePlugin) error {
	pluginPath, err := getPluginPath(plugin)
	if err != nil {
		return fmt.Errorf("failed to get plugin path: %w", err)
	}
	return updatePluginDependencies(plugin, pluginPath)
}

// getPluginPath is a helper function to get the path of a loaded plugin
func getPluginPath(plugin GitspacePlugin) (string, error) {
	pluginsDir := filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", "plugins")
	pluginPath := filepath.Join(pluginsDir, plugin.Name(), plugin.Name()+".so")

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return "", fmt.Errorf("plugin file not found: %s", pluginPath)
	}

	return pluginPath, nil
}

// Additional functions from the first attachment

func installFromGitspaceCatalog(logger *log.Logger, catalogItem string) error {
	owner := "ssotops"
	repo := "gitspace-catalog"
	defaultBranch := "master"
	catalog, err := fetchGitspaceCatalog(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to fetch Gitspace Catalog: %w", err)
	}

	plugin, ok := catalog.Plugins[catalogItem]
	if !ok {
		return fmt.Errorf("plugin %s not found in Gitspace Catalog", catalogItem)
	}

	if plugin.Path == "" {
		return fmt.Errorf("no path found for plugin %s in Gitspace Catalog", catalogItem)
	}

	rawGitHubURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, defaultBranch, plugin.Path)

	pluginsDir, err := getPluginsDir()
	if err != nil {
		return fmt.Errorf("failed to get plugins directory: %w", err)
	}

	pluginDir := filepath.Join(pluginsDir, catalogItem)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	manifestURL := fmt.Sprintf("%s/gitspace-plugin.toml", rawGitHubURL)
	manifestPath := filepath.Join(pluginDir, "gitspace-plugin.toml")
	err = downloadFile(manifestURL, manifestPath)
	if err != nil {
		return fmt.Errorf("failed to download gitspace-plugin.toml: %w", err)
	}

	soURL := fmt.Sprintf("%s/dist/%s.so", rawGitHubURL, catalogItem)
	soPath := filepath.Join(pluginDir, catalogItem+".so")
	err = downloadFile(soURL, soPath)
	if err != nil {
		return fmt.Errorf("failed to download %s.so: %w", catalogItem, err)
	}

	logger.Info("Plugin installed successfully", "name", catalogItem, "path", pluginDir)
	return nil
}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func loadPluginManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest PluginManifest
	err = toml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}
	return &manifest, nil
}

// Helper function to get plugins directory
func getPluginsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	pluginsDir := filepath.Join(homeDir, ".ssot", "gitspace", "plugins")

	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create plugins directory: %w", err)
	}

	return pluginsDir, nil
}

// This function needs to be implemented or imported from the appropriate package
func fetchGitspaceCatalog(owner, repo string) (*GitspaceCatalog, error) {
	// Implementation needed
	return nil, fmt.Errorf("fetchGitspaceCatalog not implemented")
}
