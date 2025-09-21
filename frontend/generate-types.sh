#!/usr/bin/env bash
set -euo pipefail
cd -- "$(dirname "$0")"

echo "Generating TypeScript types from todo.proto..."

# Create lib directory if it doesn't exist
mkdir -p src/lib

# Ensure protoc exists
command -v protoc >/dev/null || { echo "protoc not found. Install protoc and retry."; exit 1; }

# Use pnpm to expose protoc-gen-es from node_modules/.bin
export PATH="./node_modules/.bin:$PATH"

# Generate TypeScript types using protoc
protoc --es_out=src/lib \
  --es_opt=target=ts \
  -I ../backend \
  ../backend/todo.proto

echo "TypeScript types generated successfully!"