#!/bin/bash

# Development script with hot reloading
echo "Starting development server with hot reloading..."

# Check if Air is installed
if ! command -v air &> /dev/null; then
    echo "Installing Air for hot reloading..."
    go install github.com/air-verse/air@latest
fi

# Create tmp directory if it doesn't exist
mkdir -p tmp

# Start Air with hot reloading
air -c .air.toml