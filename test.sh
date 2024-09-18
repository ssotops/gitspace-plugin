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

# Check if gum is installed
if ! command -v gum &> /dev/null; then
    echo "gum is not installed. Please install it first:"
    echo "https://github.com/charmbracelet/gum#installation"
    exit 1
fi

# Function to run tests for a package
run_tests() {
    local package=$1
    local package_name=$2
    gum style --foreground 226 "Running tests for $package_name in directory: $(pwd)"
    gum spin --spinner dot --title "Running tests for $package_name..." -- \
        bash -c "go test -v -coverprofile=coverage.out ./$package || handle_error 'Tests failed for $package_name'"
    
    if [ ! -f coverage.out ]; then
        gum style --foreground 196 "Coverage file not generated for $package_name"
    else
        gum style \
            --foreground 82 --border-foreground 82 --border normal \
            --align left --width 70 --margin "1 2" --padding "1 2" \
            "Test Coverage for $package_name:"
        
        go tool cover -func=coverage.out
        rm coverage.out
    fi
}

# Function to run linters
run_linters() {
    local package=$1
    local package_name=$2
    gum spin --spinner dot --title "Running linters for $package_name..." -- \
        bash -c "go vet ./$package && golangci-lint run ./$package || handle_error 'Linting failed for $package_name'"
}

# ASCII Art for gitspace-plugin tester using gum
gum style \
    --foreground 212 --border-foreground 212 --border double \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Gitspace Plugin Tester"

# Summary of what's being tested
gum style \
    --foreground 226 --border-foreground 226 --border normal \
    --align left --width 70 --margin "1 2" --padding "1 2" \
    "Testing Summary:
    
1. Main gitspace-plugin package:
   - Running linters (go vet, golangci-lint)
   - Running unit tests with coverage
   - Testing core functionality (LoadPlugin, RunPlugin, etc.)

2. Example 'hello-world' plugin:
   - Running linters (go vet, golangci-lint)
   - Running unit tests with coverage
   - Testing plugin interface implementation"

# Main package tests and linting
run_linters "." "gitspace-plugin"
run_tests "." "gitspace-plugin"

# Example plugin tests and linting
gum style \
    --foreground 226 --border-foreground 226 --border normal \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "Testing Example Plugin"

cd examples/hello-world || handle_error "Failed to change to example plugin directory"
run_linters "." "hello-world plugin"
run_tests "." "hello-world plugin"
cd ../.. || handle_error "Failed to return to root directory"

gum style \
    --foreground 82 --border-foreground 82 --border double \
    --align center --width 70 --margin "1 2" --padding "1 2" \
    "All tests passed and linting successful!"
