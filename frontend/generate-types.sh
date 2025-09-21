#!/bin/bash

# Generate TypeScript types from Protocol Buffer definition
echo "Generating TypeScript types from todo.proto..."

# Create lib directory if it doesn't exist
mkdir -p src/lib

# Generate TypeScript types using protoc
protoc --es_out=src/lib \
  --es_opt=target=ts \
  --connect-es_out=src/lib \
  --connect-es_opt=target=ts \
  -I ../backend \
  ../backend/todo.proto

echo "TypeScript types generated successfully!"