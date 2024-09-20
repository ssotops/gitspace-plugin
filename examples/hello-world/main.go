//go:build linux || darwin
// +build linux darwin

package main

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/ssotops/gitspace-plugin/gsplug"
)

var Plugin HelloWorldPlugin

type HelloWorldPlugin struct {
	manifest *gsplug.PluginManifest
	logger   *log.Logger
}

func (p *HelloWorldPlugin) Init() error {
	var err error

	// Get the directory of the current executable
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)

	// Construct the path to gitspace-plugin.toml relative to the executable
	manifestPath := filepath.Join(exeDir, "gitspace-plugin.toml")

	p.manifest, err = gsplug.ReadManifest(manifestPath)
	if err != nil {
		return err
	}

	p.logger = log.New(os.Stderr)
	p.logger.SetReportCaller(true)

	return nil
}

func (p HelloWorldPlugin) Name() string {
	if p.manifest != nil {
		return p.manifest.Metadata.Name
	}
	return "hello-world"
}

func (p HelloWorldPlugin) Version() string {
	if p.manifest != nil {
		return p.manifest.Metadata.Version
	}
	return "1.0.0"
}

func (p HelloWorldPlugin) Description() string {
	if p.manifest != nil {
		return p.manifest.Metadata.Description
	}
	return "A simple Hello World plugin for Gitspace"
}

func (p HelloWorldPlugin) Run() error {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(1).
		PaddingBottom(1).
		PaddingLeft(4).
		PaddingRight(4)

	message := style.Render("Hello, World!")
	p.logger.Info(message)

	p.logger.Info("Plugin details",
		"name", p.Name(),
		"version", p.Version(),
		"description", p.Description())

	return nil
}

func (p HelloWorldPlugin) GetMenuOption() *gsplug.Option {
	key := "hello-world"
	title := "Hello World"
	if p.manifest != nil && p.manifest.Menu.Key != "" {
		key = p.manifest.Menu.Key
	}
	if p.manifest != nil && p.manifest.Menu.Title != "" {
		title = p.manifest.Menu.Title
	}
	return &gsplug.Option{
		Key:   key,
		Value: title,
	}
}

func main() {
	if err := Plugin.Init(); err != nil {
		log.Fatal("Failed to initialize plugin", "error", err)
	}

	if err := Plugin.Run(); err != nil {
		log.Fatal("Error running plugin", "error", err)
	}
}
