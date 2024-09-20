#!/bin/bash

set -e

# Function to handle errors
handle_error() {
    echo "Error: $1"
    exit 1
}

# Ensure gsplug binary exists
if [ ! -f "./cmd/gsplug/dist/gsplug" ]; then
    handle_error "gsplug binary not found. Please run build.sh first."
fi

# Ensure hello-world plugin exists
if [ ! -f "./examples/hello-world/dist/hello-world" ]; then
    handle_error "hello-world plugin not found. Please run build.sh first."
fi

# Test gsplug version
echo "Testing gsplug version..."
./cmd/gsplug/dist/gsplug version || handle_error "Failed to get gsplug version"

# Copy gitspace-plugin.toml to dist directory
echo "Copying gitspace-plugin.toml to dist directory..."
cp "./examples/hello-world/gitspace-plugin.toml" "./examples/hello-world/dist/" || handle_error "Failed to copy gitspace-plugin.toml"

# Test hello-world plugin
echo "Testing hello-world plugin..."
(cd ./examples/hello-world/dist && ./hello-world) || handle_error "Failed to run hello-world plugin"

echo "All tests passed successfully!"
