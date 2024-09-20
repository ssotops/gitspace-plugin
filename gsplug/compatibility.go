package gsplug

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
)

const (
	GitspaceVersionURL = "https://api.github.com/repos/ssotops/gitspace/releases/latest"
	VersionFile        = "gitspace-version.json"
)

// FetchLatestGitspaceVersion fetches the latest Gitspace version from GitHub
func FetchLatestGitspaceVersion() (string, error) {
	resp, err := http.Get(GitspaceVersionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	version, ok := result["tag_name"].(string)
	if !ok {
		return "", fmt.Errorf("unable to parse version from GitHub response")
	}

	return strings.TrimPrefix(version, "v"), nil
}

// UpdateVersionFile updates the local version file with the latest Gitspace version
func UpdateVersionFile() error {
	version, err := FetchLatestGitspaceVersion()
	if err != nil {
		return err
	}

	versionInfo := VersionInfo{
		GitspaceVersion: version,
		PluginAPIVersion: "1.0.0", // This should be updated manually when the plugin API changes
	}

	data, err := json.MarshalIndent(versionInfo, "", "  ")
	if err != nil {
		return err
	}

	versionFilePath := filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", VersionFile)
	return os.WriteFile(versionFilePath, data, 0644)
}

// GetVersionInfo reads the version info from the local version file
func GetVersionInfo() (*VersionInfo, error) {
	versionFilePath := filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", VersionFile)
	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		return nil, err
	}

	var versionInfo VersionInfo
	if err := json.Unmarshal(data, &versionInfo); err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

// CheckCompatibility checks if the plugin is compatible with the current Gitspace version
func CheckCompatibility(pluginVersion string) (bool, error) {
	versionInfo, err := GetVersionInfo()
	if err != nil {
		return false, err
	}

	gitspaceVersion, err := semver.NewVersion(versionInfo.GitspaceVersion)
	if err != nil {
		return false, err
	}

	pluginConstraint, err := semver.NewConstraint(pluginVersion)
	if err != nil {
		return false, err
	}

	return pluginConstraint.Check(gitspaceVersion), nil
}
