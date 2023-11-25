#!/bin/bash

# Build all platforms
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

GOOS=linux go build -o output/main $SCRIPT_DIR/demo/main.go
GOOS=darwin go build -o output/main $SCRIPT_DIR/demo/main.go
GOOS=windows go build -o output/main $SCRIPT_DIR/demo/main.go
