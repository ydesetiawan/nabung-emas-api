#!/bin/bash

# VS Code Go Setup Script
# This script ensures Go language server is properly configured

echo "🔧 Setting up Go development environment..."
echo ""

# 1. Install Go tools
echo "📦 Installing Go tools..."
go install golang.org/x/tools/gopls@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/go-delve/delve/cmd/dlv@latest

# 2. Download dependencies
echo "📥 Downloading dependencies..."
go mod download
go mod tidy

# 3. Build to verify everything works
echo "🔨 Building project..."
go build -o /tmp/nabung-emas-api cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    rm /tmp/nabung-emas-api
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "✅ Go development environment setup complete!"
echo ""
echo "📝 Next steps:"
echo "1. Reload VS Code window (Cmd+Shift+P → 'Reload Window')"
echo "2. Wait for gopls to index the project (check bottom right status bar)"
echo "3. Try Cmd+Click or F12 on any function/struct to go to definition"
echo ""
echo "💡 Tip: If go-to-definition still doesn't work:"
echo "   - Open Command Palette (Cmd+Shift+P)"
echo "   - Run 'Go: Install/Update Tools'"
echo "   - Select all tools and install"
