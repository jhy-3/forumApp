#!/bin/bash

# Install Go and Node.js without sudo

set -e

echo "========================================"
echo "  Installing Development Tools"
echo "  (Go + Node.js - No sudo required)"
echo "========================================"
echo ""

# Install Go
echo "[1/2] Installing Go..."
echo ""

if command -v go &> /dev/null; then
    echo "✓ Go is already installed: $(go version)"
else
    echo "Downloading Go 1.21.5..."
    cd /tmp
    wget -q --show-progress https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
    
    echo "Extracting..."
    tar -xzf go1.21.5.linux-amd64.tar.gz
    
    echo "Installing to ~/.local/..."
    mkdir -p ~/.local
    rm -rf ~/.local/go
    mv go ~/.local/
    
    # Add to PATH
    if ! grep -q "/.local/go/bin" ~/.bashrc; then
        echo 'export PATH=$HOME/.local/go/bin:$PATH' >> ~/.bashrc
        echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.bashrc
    fi
    
    export PATH=$HOME/.local/go/bin:$PATH
    export PATH=$HOME/go/bin:$PATH
    
    rm -f go1.21.5.linux-amd64.tar.gz
    
    echo "✓ Go installed to ~/.local/go/"
    echo "  Version: $(~/.local/go/bin/go version)"
fi

echo ""
echo "[2/2] Installing Node.js..."
echo ""

if command -v node &> /dev/null; then
    echo "✓ Node.js is already installed: $(node --version)"
else
    echo "Installing NVM (Node Version Manager)..."
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
    
    # Load NVM
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    
    echo ""
    echo "Installing Node.js 18..."
    nvm install 18
    nvm use 18
    
    echo "✓ Node.js installed"
    echo "  Version: $(node --version)"
    echo "  npm version: $(npm --version)"
fi

echo ""
echo "========================================"
echo "  Installation Complete!"
echo "========================================"
echo ""
echo "Installed tools:"
echo "  - Go: ~/.local/go/"
echo "  - Node.js: via NVM"
echo ""
echo "⚠ IMPORTANT: Reload your shell"
echo ""
echo "  Run this command:"
echo "    source ~/.bashrc"
echo ""
echo "  Or open a new terminal"
echo ""
echo "Then verify:"
echo "  go version"
echo "  node --version"
echo ""
echo "Next step:"
echo "  See: 本地运行指南.md"
echo "========================================"

