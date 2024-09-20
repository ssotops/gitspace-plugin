package gsplug

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CanonicalDeps struct {
	Versions map[string]string `json:"versions"`
}

func GetCanonicalDeps() (CanonicalDeps, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return CanonicalDeps{}, err
	}

	depsPath := filepath.Join(homeDir, ".ssot", "gitspace", "canonical-deps.json")
	data, err := ioutil.ReadFile(depsPath)
	if err != nil {
		return CanonicalDeps{}, err
	}

	var deps CanonicalDeps
	err = json.Unmarshal(data, &deps)
	if err != nil {
		return CanonicalDeps{}, err
	}

	return deps, nil
}
