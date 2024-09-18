package main

import (
	"io"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/ssotops/gitspace-plugin"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorldPlugin(t *testing.T) {
	plugin := HelloWorldPlugin{}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "hello-world", plugin.Name())
	})

	t.Run("Version", func(t *testing.T) {
		assert.Equal(t, "1.0.0", plugin.Version())
	})

	t.Run("Description", func(t *testing.T) {
		assert.Equal(t, "A simple Hello World plugin for Gitspace", plugin.Description())
	})

	t.Run("Run", func(t *testing.T) {
		logger := log.New(io.Discard)
		logger.SetLevel(log.FatalLevel)
		err := plugin.Run(logger)
		assert.NoError(t, err)
	})

	t.Run("GetMenuOption", func(t *testing.T) {
		option := plugin.GetMenuOption()
		assert.NotNil(t, option)
		assert.Equal(t, "hello-world", option.Key)
		assert.Equal(t, "Hello World", option.Value)
	})

	t.Run("Standalone", func(t *testing.T) {
		err := plugin.Standalone([]string{})
		assert.NoError(t, err)
	})

	t.Run("SetConfig", func(t *testing.T) {
		config := gitspace_plugin.PluginConfig{
			Metadata: gitspace_plugin.PluginMetadata{
				Name:    "hello-world",
				Version: "1.0.0",
			},
		}
		assert.NotPanics(t, func() {
    })
			plugin.SetConfig(config)
	})
}
