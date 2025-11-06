#!/bin/bash

# Script to fix the monorepo structure and dependencies

echo "🔧 Fixing Instashorts Monorepo Structure..."

# 1. Update Go version in go.work to match available Go version
echo "📝 Updating go.work file..."
sed -i 's/go 1.24.4/go 1.23/g' go.work

# 2. Update Go version in all go.mod files
echo "📝 Updating Go versions in go.mod files..."
sed -i 's/go 1.24.4/go 1.23/g' pkg/go.mod
sed -i 's/go 1.24.4/go 1.23/g' is-api/go.mod
sed -i 's/go 1.24.4/go 1.23/g' is-worker/go.mod

# 3. Sync workspace
echo "🔄 Syncing Go workspace..."
go work sync

# 4. Download dependencies for each module
echo "📦 Downloading dependencies..."
cd pkg && go mod tidy && cd ..
cd is-api && go mod tidy && cd ..
cd is-worker && go mod tidy && cd ..

# 5. Verify the build locally (optional)
echo "🏗️ Testing build..."
cd is-api && go build ./cmd/api && cd ..
cd is-worker && go build ./cmd/worker && cd ..

echo "✅ Monorepo structure fixed!"
echo ""
echo "Next steps:"
echo "1. Copy your vertex-ai-key.json to the root directory"
echo "2. Create a .env file with all required environment variables"
echo "3. Run 'docker-compose build' to build all services"
echo "4. Run 'docker-compose up' to start all services"