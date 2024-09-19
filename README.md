# Gitspace Plugin

![CI](https://github.com/ssotops/gitspace-plugin/workflows/CI/badge.svg)

## Overview

Gitspace Plugin is a framework for creating plugins for the Gitspace application. It provides a standardized interface for developing, testing, and integrating custom functionalities into Gitspace.

## How it works

1. Plugin Interface: The `GitspacePlugin` interface defines the methods that each plugin must implement.
2. Plugin Loading: Plugins are dynamically loaded at runtime using Go's plugin system.
3. Configuration: Plugins can be configured using TOML files.
4. Integration: Loaded plugins are integrated into Gitspace's menu system and can be executed within the Gitspace environment.

## Maintenance

### Installation

To use Gitspace Plugin in your project:

```bash
go get github.com/ssotops/gitspace-plugin
```

### Building

To build the project and its example plugins:

```bash
./build.sh
```

This script will compile the main package and the example plugins, creating both `.so` files and standalone binaries.

### Testing

To run tests for both the main package and example plugins:

```bash
./test.sh
```

This script runs linters, unit tests, and generates coverage reports.

## Development

To create a new plugin:

1. Implement the `GitspacePlugin` interface in your Go package.
2. Compile your plugin as a Go plugin (`.so` file).
3. Create a `gitspace-plugin.toml` configuration file for your plugin.
4. Place the `.so` file and `.toml` file in Gitspace's plugins directory.

For detailed examples, refer to the `examples/` directory in this repository.
