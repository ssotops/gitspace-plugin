package gsplug

type VersionInfo struct {
	GitspaceVersion string `json:"gitspace_version"`
	PluginAPIVersion string `json:"plugin_api_version"`
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
		EntryPoint string `toml:"entry_point"`
	} `toml:"sources"`
}

type Option struct {
	Key   string
	Value string
}
