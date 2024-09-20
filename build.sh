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
    cd "$dir" || handle_error "Failed to change directory to $dir"
    gum spin --spinner dot --title "Updating dependencies..." -- go get -u ./...
    go mod tidy
    cd - > /dev/null || handle_error "Failed to return to previous directory"
}

# Function to build a Go package
build_package() {
    local dir=$1
    local name=$2
    local output=$3
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
    build_output=$(go build -o "$dist_dir/$output" 2>&1)
    build_exit_code=$?

    if [ $build_exit_code -eq 0 ]; then
        gum style \
            --foreground 82 --border-foreground 82 --border normal \
            --align center --width 70 --margin "1 2" --padding "1 2" \
            "$name built successfully in $dist_dir/$output"
    else
        handle_error "Failed to build $name" "$build_output"
    fi
    
    echo "Changing back to original directory"
    cd "$original_dir" || handle_error "Failed to return to original directory"
    echo "Current working directory after returning: $(pwd)"
}

# Update main gitspace-plugin dependencies
update_dependencies .

# Build gsplug package
gum spin --spinner dot --title "Building gsplug package..." -- go build ./gsplug || handle_error "Failed to build gsplug package"

# Build cmd/gsplug
echo "Current working directory before building cmd/gsplug: $(pwd)"
if [ -d "cmd/gsplug" ]; then
    echo "cmd/gsplug directory exists"
    build_package "cmd/gsplug" "gsplug CLI" "gsplug"
else
    echo "cmd/gsplug directory does not exist"
    gum style \
        --foreground 208 --border-foreground 208 --border normal \
        --align center --width 70 --margin "1 2" --padding "1 2" \
        "cmd/gsplug directory does not exist. Skipping build of gsplug CLI."
fi

# Build examples/hello-world
build_package "examples/hello-world" "hello-world plugin" "hello-world.so"

# Also build hello-world as a standalone binary
build_package "examples/hello-world" "hello-world standalone" "hello-world"

# Print summary
gum style \
    --foreground 226 --border-foreground 226 --border double \
    --align left --width 70 --margin "1 2" --padding "1 2" \
    "Build Summary:

1. gsplug package: Built ✅
2. gsplug CLI: Built ✅
3. hello-world plugin: Built ✅
4. hello-world standalone: Built ✅

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
