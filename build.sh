#!/bin/bash

set -o pipefail

# Function to handle errors
handle_error() {
    gum style \
        --foreground 196 --border-foreground 196 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "An error occurred: $1"
    exit 1
}

# Function to install gum
install_gum() {
    echo "Installing gum..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install gum || handle_error "Failed to install gum"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "For Ubuntu/Debian:"
        echo "sudo mkdir -p /etc/apt/keyrings"
        echo "curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg"
        echo 'echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list'
        echo "sudo apt update && sudo apt install gum"
        echo ""
        echo "For other Linux distributions, please visit: https://github.com/charmbracelet/gum#installation"
        read -p "Do you want to proceed with the installation for Ubuntu/Debian? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            sudo mkdir -p /etc/apt/keyrings || handle_error "Failed to create keyrings directory"
            curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg || handle_error "Failed to download and save GPG key"
            echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list || handle_error "Failed to add Charm repository"
            sudo apt update && sudo apt install gum || handle_error "Failed to install gum"
        else
            echo "Please install gum manually and run this script again."
            exit 1
        fi
    else
        echo "Unsupported operating system. Please install gum manually:"
        echo "https://github.com/charmbracelet/gum#installation"
        exit 1
    fi
}

# Check if gum is installed
if ! command -v gum &> /dev/null; then
    echo "gum is not installed."
    read -p "Do you want to install gum? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        install_gum
    else
        echo "Please install gum manually and run this script again."
        exit 1
    fi
fi

# Function to update Charm library versions
update_charm_versions() {
    local dir=$1
    cd "$dir" || handle_error "Failed to change directory to $dir"
    go get github.com/charmbracelet/huh@latest || handle_error "Failed to update huh"
    go get github.com/charmbracelet/log@latest || handle_error "Failed to update log"
    go mod tidy || handle_error "Failed to tidy go.mod"
    cd - > /dev/null || handle_error "Failed to return to previous directory"
}

# ASCII Art for gitspace_plugin builder using gum
gum style \
    --foreground 212 --border-foreground 212 --border double \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Gitspace Plugin Builder"

# Update main package
gum spin --spinner dot --title "Updating main package..." -- bash -c "update_charm_versions . || handle_error 'Failed to update main package'"

# Build main package
gum spin --spinner dot --title "Building main package..." -- bash -c "
    mkdir -p dist || handle_error 'Failed to create dist directory'
    go build -o dist/gitspace-plugin . || handle_error 'Failed to build main package'
"

# Run tests
gum spin --spinner dot --title "Running tests..." -- bash -c "go test ./... || handle_error 'Some tests failed'"

# Build example plugin
gum spin --spinner dot --title "Building example plugin..." -- bash -c "
    cd examples/hello-world || handle_error 'Failed to change to example plugin directory'
    update_charm_versions . || handle_error 'Failed to update example plugin'
    go mod edit -replace=github.com/ssotops/gitspace-plugin=../../ || handle_error 'Failed to edit go.mod'
    go mod tidy || handle_error 'Failed to tidy example plugin go.mod'
    mkdir -p dist || handle_error 'Failed to create dist directory'
    CGO_ENABLED=1 go build -buildmode=plugin -o dist/hello-world.so . || handle_error 'Failed to build example plugin .so'
    go build -o dist/hello-world . || handle_error 'Failed to build example plugin standalone binary'
    cd ../.. || handle_error 'Failed to return to root directory'
"

# Verify builds
if [ ! -f "dist/gitspace-plugin" ]; then
    handle_error "Main gitspace-plugin binary not found"
fi

if [ ! -f "examples/hello-world/dist/hello-world.so" ]; then
    handle_error "Example plugin .so file not found"
fi

if [ ! -f "examples/hello-world/dist/hello-world" ]; then
    handle_error "Example plugin standalone binary not found"
fi

gum style \
    --foreground 212 --border-foreground 212 --border normal \
    --align left --width 70 --margin "1 2" --padding "1 2" \
    "Build complete!
Gitspace Plugin package: ./dist/gitspace-plugin
Example plugin .so: ./examples/hello-world/dist/hello-world.so
Example plugin binary: ./examples/hello-world/dist/hello-world"

# Ask if user wants to install the example plugin
if gum confirm "Do you want to install the example plugin to ~/.gitspace/plugins?"; then
    # Create plugins directory if it doesn't exist
    mkdir -p ~/.gitspace/plugins || handle_error "Failed to create plugins directory"

    # Copy example plugin to the plugins directory
    cp examples/hello-world/dist/hello-world.so ~/.gitspace/plugins/ || handle_error "Failed to copy plugin .so file"
    cp examples/hello-world/gitspace-plugin.toml ~/.gitspace/plugins/hello-world.toml || handle_error "Failed to copy plugin toml file"

    gum style \
        --foreground 82 --border-foreground 82 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "Example plugin installed to ~/.gitspace/plugins/"
else
    gum style \
        --foreground 208 --border-foreground 208 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "Example plugin was not installed."
fi

# Print installed plugins
echo "Currently installed plugins:"
for plugin in ~/.gitspace/plugins/*.so; do
    if [ -f "$plugin" ]; then
        plugin_name=$(basename "$plugin" .so)
        gum style \
            --foreground 39 --border-foreground 39 --border normal \
            --align left --width 50 --margin "0 2" --padding "0 1" \
            "🔌 $plugin_name"
    fi
done

# Print tree structure of plugins directory
if command -v tree &> /dev/null; then
    tree_output=$(tree -L 2 ~/.gitspace/plugins)
    gum style \
        --foreground 226 --border-foreground 226 --border double \
        --align left --width 70 --margin "1 2" --padding "1 2" \
        "Plugins Directory Structure:

$tree_output"
else
    gum style \
        --foreground 226 --border-foreground 226 --border double \
        --align left --width 70 --margin "1 2" --padding "1 2" \
        "Plugins Directory Structure:

$(ls -R ~/.gitspace/plugins)"
fi

gum style \
    --foreground 214 --border-foreground 214 --border normal \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Note: This build script is for the gitspace_plugin package. 
    To use this package in your Gitspace project, make sure to update 
    your Gitspace build script and import paths accordingly."
