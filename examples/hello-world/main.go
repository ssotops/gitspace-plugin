package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	gitspace_plugin "github.com/ssotops/gitspace-plugin"
)

type HelloWorldPlugin struct{}

var Plugin HelloWorldPlugin

func (p HelloWorldPlugin) Name() string {
	return "hello-world"
}

func (p HelloWorldPlugin) Version() string {
	return "1.0.0"
}

func (p HelloWorldPlugin) Description() string {
	return "A simple Hello World plugin for Gitspace"
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

func (p HelloWorldPlugin) SetConfig(config gitspace_plugin.PluginConfig) {
	// This plugin doesn't use any configuration, but we need to implement this method
}

func main() {
	plugin := HelloWorldPlugin{}
	if err := plugin.Standalone(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
