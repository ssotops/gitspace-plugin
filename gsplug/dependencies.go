package gsplug

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const GitspaceRepoURL = "https://raw.githubusercontent.com/ssotops/gitspace/main/go.mod"

// EnsureGitspaceModFile checks if the gitspace-go.mod file exists, and if not, downloads it
func EnsureGitspaceModFile() error {
	gitspaceModPath := filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", "plugins", "gitspace-go.mod")

	if _, err := os.Stat(gitspaceModPath); os.IsNotExist(err) {
		// Download the file
		if err := downloadGitspaceMod(gitspaceModPath); err != nil {
			return fmt.Errorf("failed to download gitspace-go.mod: %w\nPlease manually add the file from %s", err, GitspaceRepoURL)
		}
	}

	return nil
}

// downloadGitspaceMod downloads the go.mod file from the Gitspace repository
func downloadGitspaceMod(destPath string) error {
	resp, err := http.Get(GitspaceRepoURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// GetGitspaceDependencies parses the gitspace-go.mod file and returns a map of dependencies
func GetGitspaceDependencies() (map[string]string, error) {
	gitspaceModPath := filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", "plugins", "gitspace-go.mod")

	if err := EnsureGitspaceModFile(); err != nil {
		return nil, err
	}

	return parseDependencies(gitspaceModPath)
}

// parseDependencies reads a go.mod file and returns a map of dependencies
func parseDependencies(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	deps := make(map[string]string)
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^(\S+)\s+(\S+)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "module") || strings.HasPrefix(line, "go ") {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			deps[matches[1]] = matches[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return deps, nil
}

// MergeDependencies combines plugin dependencies with Gitspace dependencies, preferring Gitspace versions
func MergeDependencies(pluginDeps, gitspaceDeps map[string]string) map[string]string {
	mergedDeps := make(map[string]string)

	for dep, version := range pluginDeps {
		if gitspaceVersion, exists := gitspaceDeps[dep]; exists {
			mergedDeps[dep] = gitspaceVersion
		} else {
			mergedDeps[dep] = version
		}
	}

	for dep, version := range gitspaceDeps {
		if _, exists := mergedDeps[dep]; !exists {
			mergedDeps[dep] = version
		}
	}

	return mergedDeps
}

// UpdatePluginDependencies updates the plugin's go.mod file with merged dependencies
func UpdatePluginDependencies(pluginDir string) error {
	pluginModPath := filepath.Join(pluginDir, "go.mod")
	pluginDeps, err := parseDependencies(pluginModPath)
	if err != nil {
		return fmt.Errorf("failed to parse plugin dependencies: %w", err)
	}

	gitspaceDeps, err := GetGitspaceDependencies()
	if err != nil {
		return fmt.Errorf("failed to get Gitspace dependencies: %w", err)
	}

	mergedDeps := MergeDependencies(pluginDeps, gitspaceDeps)

	// Write updated go.mod
	return writeGoMod(pluginModPath, mergedDeps)
}

// writeGoMod writes the merged dependencies to the plugin's go.mod file
func writeGoMod(path string, deps map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write module declaration (preserving the original module name)
	originalModule, err := getModuleName(path)
	if err != nil {
		return fmt.Errorf("failed to get original module name: %w", err)
	}
	_, err = fmt.Fprintf(writer, "module %s\n\n", originalModule)
	if err != nil {
		return err
	}

	// Write go version (using the same version as Gitspace)
	gitspaceGoVersion, err := getGoVersion(filepath.Join(os.Getenv("HOME"), ".ssot", "gitspace", "plugins", "gitspace-go.mod"))
	if err != nil {
		return fmt.Errorf("failed to get Gitspace Go version: %w", err)
	}
	_, err = fmt.Fprintf(writer, "go %s\n\n", gitspaceGoVersion)
	if err != nil {
		return err
	}

	// Write dependencies
	_, err = fmt.Fprintln(writer, "require (")
	if err != nil {
		return err
	}
	for dep, version := range deps {
		_, err = fmt.Fprintf(writer, "\t%s %s\n", dep, version)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(writer, ")")
	if err != nil {
		return err
	}

	return nil
}

// getModuleName retrieves the module name from a go.mod file
func getModuleName(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module declaration not found in go.mod")
}

// getGoVersion retrieves the Go version from a go.mod file
func getGoVersion(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "go ") {
			return strings.TrimPrefix(line, "go "), nil
		}
	}

	return "", fmt.Errorf("go version not found in go.mod")
}
