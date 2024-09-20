package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ssotops/gitspace-plugin/gsplug"
)

func main() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildAll := buildCmd.Bool("all", false, "Build all plugins")

	updateDepsCmd := flag.NewFlagSet("update-deps", flag.ExitOnError)

	updateVersionCmd := flag.NewFlagSet("update-version", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'build', 'update-deps', or 'update-version' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		buildCmd.Parse(os.Args[2:])
		if *buildAll {
			if err := gsplug.BuildAllPlugins(); err != nil {
				fmt.Printf("Error building all plugins: %v\n", err)
				os.Exit(1)
			}
		} else {
			if buildCmd.NArg() < 1 {
				fmt.Println("Please specify a plugin directory")
				os.Exit(1)
			}
			pluginDir := buildCmd.Arg(0)
			if err := gsplug.BuildPlugin(pluginDir); err != nil {
				fmt.Printf("Error building plugin: %v\n", err)
				os.Exit(1)
			}
		}

	case "update-deps":
		updateDepsCmd.Parse(os.Args[2:])
		if updateDepsCmd.NArg() < 1 {
			fmt.Println("Please specify a plugin directory")
			os.Exit(1)
		}
		pluginDir := updateDepsCmd.Arg(0)
		if err := gsplug.UpdatePluginDependencies(pluginDir); err != nil {
			fmt.Printf("Error updating plugin dependencies: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Plugin dependencies updated successfully")

	case "update-version":
		updateVersionCmd.Parse(os.Args[2:])
		if err := gsplug.UpdateVersionFile(); err != nil {
			fmt.Printf("Error updating version file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Version file updated successfully")

	default:
		fmt.Println("Expected 'build', 'update-deps', or 'update-version' subcommands")
		os.Exit(1)
	}
}
