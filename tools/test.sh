#!/bin/bash

set -e

MODE="$1"
DEFAULT_FLAGS="-covermode=atomic -coverpkg=./... -coverprofile=coverage.out"

case "$MODE" in
  race)
    echo "🔬 Running tests with race detector..."
    go test ./... $DEFAULT_FLAGS -race
    ;;
  verbose)
    echo "🗣️ Running verbose tests..."
    go test ./... $DEFAULT_FLAGS -v
    ;;
  coverage)
    echo "📊 Opening coverage report..."
    go tool cover -html=coverage.out
    ;;
  *)
    echo "🧪 Running standard tests..."
    go test ./... $DEFAULT_FLAGS
    ;;
esac
