#!/bin/bash
set -e

echo "Generating protobuf files..."

# Ensure the output directory exists
mkdir -p proto/plugin/v1

# Generate Go code from proto files
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/plugin/v1/plugin.proto

echo "Protobuf generation complete!"