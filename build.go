package gsplug

import (
	"os"
	"os/exec"
	"path/filepath"
)

func BuildPlugin(pluginDir string) error {
	// Ensure gitspace dependencies are used
	if err := updateDependencies(pluginDir); err != nil {
		return err
	}

	// Build the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", filepath.Join(pluginDir, "plugin.so"))
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}


