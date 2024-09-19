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

# ASCII Art for gitspace-plugin builder using gum
gum style \
    --foreground 212 --border-foreground 212 --border double \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Gitspace Plugin Builder"

# Update main package
gum spin --spinner dot --title "Updating main package..." -- bash -c "update_charm_versions . || handle_error 'Failed to update main package'"

# Run tests for main package
gum spin --spinner dot --title "Running tests for main package..." -- bash -c "go test -v ./... 2>&1 || handle_error 'Some tests failed'"

# Define the correct plugin installation directory
PLUGIN_DIR="$HOME/.ssot/gitspace/plugins/hello-world"

# Build example plugin
echo "Building example plugin..."
(
    cd examples/hello-world || handle_error 'Failed to change to example plugin directory'
    
    if [ ! -f go.mod ]; then
        go mod init github.com/ssotops/gitspace-plugin/examples/hello-world || handle_error 'Failed to initialize Go module'
    else
        go mod edit -module=github.com/ssotops/gitspace-plugin/examples/hello-world || handle_error 'Failed to update module name'
    fi
    
    go mod edit -replace=github.com/ssotops/gitspace-plugin=../../ || handle_error 'Failed to edit go.mod'

    go get github.com/ssotops/gitspace-plugin@latest || handle_error 'Failed to update local gitspace-plugin'
    go mod tidy || handle_error 'Failed to tidy go.mod after updating gitspace-plugin'
    
    go get github.com/charmbracelet/huh@latest || handle_error 'Failed to update huh'
    go get github.com/charmbracelet/log@latest || handle_error 'Failed to update log'
    go get github.com/ssotops/gitspace-plugin@latest || handle_error 'Failed to update gitspace-plugin'
    
    go mod tidy || handle_error 'Failed to tidy example plugin go.mod'
    
    mkdir -p dist || handle_error 'Failed to create dist directory'
    
    CGO_ENABLED=1 go build -v -buildmode=plugin -o dist/hello-world.so . || handle_error 'Failed to build example plugin .so'

    go build -v -o dist/hello-world . || handle_error 'Failed to build example plugin standalone binary'
    
    echo "Contents of dist directory:"
    ls -l dist

    cd ../.. || handle_error 'Failed to return to root directory'
) 2>&1 | tee build.log

# Check if the build was successful
if [ ! -f "examples/hello-world/dist/hello-world.so" ] || [ ! -f "examples/hello-world/dist/hello-world" ]; then
    echo "Build failed. See build.log for details."
    cat build.log
    handle_error "Build failed"
fi

echo "Build verification completed successfully."

gum style \
    --foreground 212 --border-foreground 212 --border normal \
    --align left --width 70 --margin "1 2" --padding "1 2" \
    "Build complete!
Example plugin .so: ./examples/hello-world/dist/hello-world.so
Example plugin binary: ./examples/hello-world/dist/hello-world"

# Ask if user wants to install the example plugin
if gum confirm "Do you want to install the example plugin to $PLUGIN_DIR?"; then
    # Create plugins directory if it doesn't exist
    mkdir -p "$PLUGIN_DIR" || handle_error "Failed to create plugins directory"

    # Copy the built .so file
    cp examples/hello-world/dist/hello-world.so "$PLUGIN_DIR/" || handle_error "Failed to copy plugin .so file"
    
    # Optionally copy the binary (uncomment if needed)
    # cp examples/hello-world/dist/hello-world "$PLUGIN_DIR/" || handle_error "Failed to copy plugin binary"

    # Copy configuration files only if they don't exist
    if [ ! -f "$PLUGIN_DIR/gitspace-plugin.toml" ]; then
        cp examples/hello-world/gitspace-plugin.toml "$PLUGIN_DIR/" || handle_error "Failed to copy plugin toml file"
    fi
    if [ ! -f "$PLUGIN_DIR/go.mod" ]; then
        cp examples/hello-world/go.mod "$PLUGIN_DIR/" || handle_error "Failed to copy go.mod file"
    fi

    gum style \
        --foreground 82 --border-foreground 82 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "Example plugin installed to $PLUGIN_DIR"
else
    gum style \
        --foreground 208 --border-foreground 208 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "Example plugin was not installed."
fi


# Print installed plugins
echo "Currently installed plugins:"
for plugin in "$PLUGIN_DIR"/*.so; do
    if [ -f "$plugin" ]; then
        plugin_name=$(basename "$plugin" .so)
        gum style \
            --foreground 39 --border-foreground 39 --border normal \
            --align left --width 50 --margin "0 2" --padding "0 1" \
            "🔌 $plugin_name"
    fi
done

# Update the tree output to use the correct directory
if command -v tree &> /dev/null; then
    tree_output=$(tree -L 2 "$PLUGIN_DIR")
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

$(ls -R "$PLUGIN_DIR")"
fi

gum style \
    --foreground 214 --border-foreground 214 --border normal \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Note: The gitspace-plugin package is built as a library for other packages to import.
    No local binary is produced for the main package."
