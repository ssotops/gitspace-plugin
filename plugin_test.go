package gitspace_plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

type MockPlugin struct {
	name        string
	version     string
	description string
	config      PluginConfig
}

func (m *MockPlugin) Name() string                       { return m.name }
func (m *MockPlugin) Version() string                    { return m.version }
func (m *MockPlugin) Description() string                { return m.description }
func (m *MockPlugin) Run(logger *log.Logger) error       { return nil }
func (m *MockPlugin) GetMenuOption() *huh.Option[string] { return nil }
func (m *MockPlugin) Standalone(args []string) error     { return nil }
func (m *MockPlugin) SetConfig(config PluginConfig)      { m.config = config }

func TestLoadPluginWithConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "plugin-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock gitspace-plugin.toml file
	configContent := `
[metadata]
name = "mock-plugin"
version = "1.0.0"
description = "A mock plugin for testing"
author = "Test Author"
tags = ["test", "mock"]

[menu]
title = "Mock Plugin"
key = "mock-plugin"
`
	err = os.WriteFile(filepath.Join(tempDir, "gitspace-plugin.toml"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Mock the LoadPlugin function
	originalLoadPlugin := LoadPlugin
	defer func() { LoadPlugin = originalLoadPlugin }()
	LoadPlugin = func(pluginPath string) (GitspacePlugin, error) {
		return &MockPlugin{
			name:        "mock-plugin",
			version:     "1.0.0",
			description: "A mock plugin for testing",
		}, nil
	}

	// Test LoadPluginWithConfig
	plugin, err := LoadPluginWithConfig(filepath.Join(tempDir, "mock-plugin.so"))
	assert.NoError(t, err)
	assert.NotNil(t, plugin)

	mockPlugin, ok := plugin.(*MockPlugin)
	assert.True(t, ok)
	assert.Equal(t, "mock-plugin", mockPlugin.Name())
	assert.Equal(t, "1.0.0", mockPlugin.Version())
	assert.Equal(t, "A mock plugin for testing", mockPlugin.Description())
	assert.Equal(t, "Test Author", mockPlugin.config.Metadata.Author)
}

func TestRunPlugin(t *testing.T) {
	mockPlugin := &MockPlugin{
		name:        "mock-plugin",
		version:     "1.0.0",
		description: "A mock plugin for testing",
	}

	logger := log.New(os.Stderr)
	err := RunPlugin(mockPlugin, logger)
	assert.NoError(t, err)
}

func TestRunStandalonePlugin(t *testing.T) {
	mockPlugin := &MockPlugin{
		name:        "mock-plugin",
		version:     "1.0.0",
		description: "A mock plugin for testing",
	}

	err := RunStandalonePlugin(mockPlugin, []string{"arg1", "arg2"})
	assert.NoError(t, err)
}
