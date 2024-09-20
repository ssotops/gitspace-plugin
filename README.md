# gsplug CLI Tool

`gsplug` is a command-line interface (CLI) tool for managing Gitspace plugins. It provides functionality for building plugins and updating their dependencies.

## Installation

To install `gsplug`, follow these steps:

1. Clone the gitspace-plugin repository:
   ```
   git clone https://github.com/ssotops/gitspace-plugin.git
   ```

2. Navigate to the repository directory:
   ```
   cd gitspace-plugin
   ```

3. Build the `gsplug` tool:
   ```
   go build -o gsplug cmd/gsplug/main.go
   ```

4. (Optional) Move the `gsplug` binary to a directory in your PATH for easier access:
   ```
   sudo mv gsplug /usr/local/bin/
   ```

## Usage

### Building Plugins

To build a single plugin:
```
gsplug build /path/to/plugin
```

To build all plugins in the Gitspace plugins directory:
```
gsplug build -all
```

### Updating Plugin Dependencies

To update the dependencies of a plugin:
```
gsplug update-deps /path/to/plugin
```

## Examples

1. Build a specific plugin:
   ```
   gsplug build ~/.ssot/gitspace/plugins/my-plugin
   ```

2. Build all plugins:
   ```
   gsplug build -all
   ```

3. Update dependencies for a specific plugin:
   ```
   gsplug update-deps ~/.ssot/gitspace/plugins/my-plugin
   ```

## Note

Ensure that you have the necessary permissions to access and modify the plugin directories. The tool assumes that plugins are located in the `~/.ssot/gitspace/plugins/` directory by default.

For any issues or feature requests, please open an issue in the [gitspace-plugin repository](https://github.com/ssotops/gitspace-plugin).
