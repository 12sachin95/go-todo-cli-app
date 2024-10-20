#!/bin/bash

# Name of the CLI binary
CLI_NAME="todo-cli"

# Step 1: Build the Go project
echo "Building the CLI project..."
go build -o $CLI_NAME

# Check if build succeeded
if [ $? -ne 0 ]; then
  echo "Build failed!"
  exit 1
fi

echo "Build succeeded! Running the CLI app..."

# Step 2: Run the built CLI application with the provided arguments
./$CLI_NAME "$@"
