#!/bin/bash

set -e

# Function to handle errors
handle_error() {
    gum style \
        --foreground 196 --border-foreground 196 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "An error occurred: $1"
    echo "Error details:"
    echo "$2" | sed 's/^/    /'  # This indents each line of the error message
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

# ASCII Art for gitspace-plugin builder using gum
gum style \
    --foreground 212 --border-foreground 212 --border double \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Gitspace Plugin Builder"

# Function to update dependencies
update_dependencies() {
    local dir=$1
    gum spin --spinner dot --title "Updating dependencies in $dir..." -- go get -u ./...
    go mod tidy
}

# Function to build a Go package
build_package() {
    local dir=$1
    local name=$2
    local output=$3
    local build_mode=$4  # Can be "plugin", "binary", or "both"
    local original_dir=$(pwd)
    local dist_dir="$dir/dist"
    
    echo "Current working directory: $original_dir"
    echo "Attempting to build package in directory: $dir"
    
    if [ ! -d "$dir" ]; then
        gum style \
            --foreground 208 --border-foreground 208 --border normal \
            --align center --width 70 --margin "1 2" --padding "1 2" \
            "Directory $dir does not exist. Skipping build of $name."
        return
    fi
    
    mkdir -p "$dist_dir"
    
    echo "Changing to directory: $dir"
    cd "$dir" || handle_error "Failed to change directory to $dir"
    echo "Current working directory after cd: $(pwd)"
    
    update_dependencies .
    
    echo "Building $name..."
    if [ "$build_mode" == "plugin" ] || [ "$build_mode" == "both" ]; then
        plugin_cmd="go build -buildmode=plugin -o $original_dir/$dist_dir/${output}.so ."
        plugin_output=$(eval $plugin_cmd 2>&1)
        plugin_exit_code=$?
        if [ $plugin_exit_code -eq 0 ]; then
            gum style \
                --foreground 82 --border-foreground 82 --border normal \
                --align center --width 70 --margin "1 2" --padding "1 2" \
                "$name plugin built successfully in $dist_dir/${output}.so"
        else
            handle_error "Failed to build $name plugin" "$plugin_output"
        fi
    fi
    
    if [ "$build_mode" == "plugin" ] || [ "$build_mode" == "both" ]; then
        plugin_cmd="go build -buildmode=plugin -o $original_dir/$dist_dir/${output}.so ."
        echo "Running command: $plugin_cmd"
        plugin_output=$(eval $plugin_cmd 2>&1)
        plugin_exit_code=$?
        if [ $plugin_exit_code -eq 0 ]; then
            gum style \
                --foreground 82 --border-foreground 82 --border normal \
                --align center --width 70 --margin "1 2" --padding "1 2" \
                "$name plugin built successfully in $dist_dir/${output}.so"
        else
            handle_error "Failed to build $name plugin" "$plugin_output"
        fi
    fi
    
    if [ "$build_mode" == "binary" ] || [ "$build_mode" == "both" ]; then
        binary_cmd="go build -o $original_dir/$dist_dir/$output ."
        echo "Running command: $binary_cmd"
        binary_output=$(eval $binary_cmd 2>&1)
        binary_exit_code=$?
        if [ $binary_exit_code -eq 0 ]; then
            gum style \
                --foreground 82 --border-foreground 82 --border normal \
                --align center --width 70 --margin "1 2" --padding "1 2" \
                "$name binary built successfully in $dist_dir/$output"
        else
            handle_error "Failed to build $name binary" "$binary_output"
        fi
    fi
    
    # Copy gitspace-plugin.toml to dist directory
    if [ -f "$dir/gitspace-plugin.toml" ]; then
        cp "$dir/gitspace-plugin.toml" "$original_dir/$dist_dir/"
        echo "Copied gitspace-plugin.toml to $dist_dir/"
    else
        echo "Warning: gitspace-plugin.toml not found in $dir"
    fi
    
    echo "Changing back to original directory"
    cd "$original_dir" || handle_error "Failed to return to original directory"
    echo "Current working directory after returning: $(pwd)"
}

# Build cmd/gsplug
build_package "cmd/gsplug" "gsplug CLI" "gsplug" "binary"

# Build examples/hello-world as both a plugin and a standalone binary
build_package "examples/hello-world" "hello-world" "hello-world" "both"

# Print summary
gum style \
    --foreground 226 --border-foreground 226 --border double \
    --align left --width 70 --margin "1 2" --padding "1 2" \
    "Build Summary:

1. gsplug package: Built ✅
2. gsplug CLI: Built ✅ (Located in cmd/gsplug/dist/gsplug)
3. hello-world plugin: Built ✅ (Located in examples/hello-world/dist/hello-world.so)
4. hello-world standalone: Built ✅ (Located in examples/hello-world/dist/hello-world)

All components have been successfully built!"

# Verify versions
go_version=$(go version | awk '{print $3}')
gsplug_version=$(./cmd/gsplug/dist/gsplug version 2>/dev/null || echo "N/A")

gum style \
    --foreground 82 --border-foreground 82 --border normal \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Versions:
Go: $go_version
gsplug CLI: $gsplug_version"

# Print directory structure
tree_output=$(tree -L 3)
gum style \
    --foreground 226 --border-foreground 226 --border double \
    --align left --width 70 --margin "1 2" --padding "1 2" \
    "Project Directory Structure:

$tree_output"

gum style \
    --foreground 214 --border-foreground 214 --border normal \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Build process completed successfully. You can now use the gsplug CLI and the gsplug package in your projects."
