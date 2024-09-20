package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/ssotops/gitspace-plugin"
)

type HelloWorldPlugin struct {
	config gitspace_plugin.PluginConfig
}

var Plugin = &HelloWorldPlugin{}

func (p HelloWorldPlugin) Name() string {
	return p.config.Metadata.Name
}

func (p HelloWorldPlugin) Version() string {
	return p.config.Metadata.Version
}

func (p HelloWorldPlugin) Description() string {
	return p.config.Metadata.Description
}

func (p HelloWorldPlugin) Run(logger *log.Logger) error {
	logger.Info("Hello from the Hello World plugin!")
	return nil
}

func (p HelloWorldPlugin) GetMenuOption() *gitspace_plugin.Option {
	return &gitspace_plugin.Option{
		Key:   "hello-world",
		Value: "Hello World",
	}
}

func (p HelloWorldPlugin) Standalone(args []string) error {
	fmt.Println("Hello from the standalone Hello World plugin!")
	return nil
}

func (p *HelloWorldPlugin) SetConfig(config gitspace_plugin.PluginConfig) {
	p.config = config
}

func (p HelloWorldPlugin) GetDependencies() map[string]string {
	return gitspace_plugin.GetSharedDependencies()
}

func (p *HelloWorldPlugin) Update(newPlugin gitspace_plugin.gsplug) {
	if newP, ok := newPlugin.(*HelloWorldPlugin); ok {
		*p = *newP
	}
}

func main() {
	plugin := HelloWorldPlugin{}
	if err := plugin.Standalone(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
