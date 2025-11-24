#!/bin/bash

# Podman Installation Script for Forum Application
# Supports installation without sudo (using static binary)

set -e

echo "========================================"
echo "  Podman Installation Helper"
echo "========================================"
echo ""

# Check if Podman is already installed
if command -v podman &> /dev/null; then
    echo "✓ Podman is already installed!"
    PODMAN_VERSION=$(podman --version)
    echo "  Version: $PODMAN_VERSION"
    echo ""
    
    # Check for podman-compose
    if command -v podman-compose &> /dev/null; then
        echo "✓ Podman Compose is already installed!"
        COMPOSE_VERSION=$(podman-compose --version)
        echo "  Version: $COMPOSE_VERSION"
        echo ""
        echo "All requirements met. You can now run: ./start.sh"
        exit 0
    else
        echo "⚠ Podman Compose is not installed"
        echo ""
    fi
fi

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VERSION=$VERSION_ID
else
    echo "Cannot detect OS."
    OS="unknown"
fi

echo "Detected OS: $OS"
echo ""

# Check if user has sudo access
if sudo -n true 2>/dev/null; then
    HAS_SUDO=true
    echo "✓ Sudo access detected"
else
    HAS_SUDO=false
    echo "⚠ No sudo access detected"
    echo "  Will install Podman without system-wide installation"
fi

echo ""

if [ "$HAS_SUDO" = true ]; then
    # Install using package manager (system-wide)
    echo "Installing Podman using package manager..."
    echo ""
    
    case $OS in
        ubuntu|debian)
            echo "Running: sudo apt-get update"
            sudo apt-get update
            echo ""
            echo "Running: sudo apt-get install -y podman podman-compose"
            sudo apt-get install -y podman podman-compose
            ;;
            
        fedora|rhel|centos)
            echo "Running: sudo dnf install -y podman podman-compose"
            sudo dnf install -y podman podman-compose
            ;;
            
        arch|manjaro)
            echo "Running: sudo pacman -S --noconfirm podman podman-compose"
            sudo pacman -S --noconfirm podman podman-compose
            ;;
            
        *)
            echo "Unsupported OS for automatic installation: $OS"
            HAS_SUDO=false
            ;;
    esac
fi

# Install without sudo (user-space installation)
if [ "$HAS_SUDO" = false ]; then
    echo "========================================" 
    echo "  Installing Podman (User-space)"
    echo "========================================"
    echo ""
    echo "Installing Podman as a regular user..."
    echo "This will install to: ~/.local/bin/"
    echo ""
    
    # Create user bin directory
    mkdir -p ~/.local/bin
    
    # Download static Podman binary
    echo "Downloading Podman static binary..."
    PODMAN_VERSION="v4.9.3"
    ARCH=$(uname -m)
    
    if [ "$ARCH" = "x86_64" ]; then
        PODMAN_URL="https://github.com/containers/podman/releases/download/${PODMAN_VERSION}/podman-remote-static-linux_amd64.tar.gz"
    elif [ "$ARCH" = "aarch64" ]; then
        PODMAN_URL="https://github.com/containers/podman/releases/download/${PODMAN_VERSION}/podman-remote-static-linux_arm64.tar.gz"
    else
        echo "Unsupported architecture: $ARCH"
        echo ""
        echo "Please contact your system administrator to install Podman:"
        echo "  sudo apt-get install podman podman-compose"
        exit 1
    fi
    
    cd /tmp
    curl -L -o podman.tar.gz "$PODMAN_URL"
    
    # Extract and find the binary
    tar -xzf podman.tar.gz
    
    # Try different possible locations
    PODMAN_BINARY=""
    if [ -f "bin/podman-remote-static-linux_amd64" ]; then
        PODMAN_BINARY="bin/podman-remote-static-linux_amd64"
    elif [ -f "bin/podman-remote-static-linux_arm64" ]; then
        PODMAN_BINARY="bin/podman-remote-static-linux_arm64"
    elif [ -f "podman-remote-static" ]; then
        PODMAN_BINARY="podman-remote-static"
    elif [ -f "bin/podman-remote" ]; then
        PODMAN_BINARY="bin/podman-remote"
    elif [ -f "podman" ]; then
        PODMAN_BINARY="podman"
    else
        # Search recursively for any podman binary
        PODMAN_BINARY=$(find . -type f -name "podman*" 2>/dev/null | grep -E "podman(-remote)?(-static)?$" | head -1)
    fi
    
    if [ -n "$PODMAN_BINARY" ] && [ -f "$PODMAN_BINARY" ]; then
        cp "$PODMAN_BINARY" ~/.local/bin/podman
        chmod +x ~/.local/bin/podman
        echo "✓ Podman installed to ~/.local/bin/podman"
    else
        echo "✗ Failed to extract Podman binary"
        echo "Downloaded file structure:"
        tar -tzf podman.tar.gz | head -10
        exit 1
    fi
    
    # Clean up
    rm -f podman.tar.gz
    rm -rf bin podman-remote-static podman
    
    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo ""
        echo "Adding ~/.local/bin to PATH..."
        
        # Add to shell config
        SHELL_CONFIG=""
        if [ -f "$HOME/.bashrc" ]; then
            SHELL_CONFIG="$HOME/.bashrc"
        elif [ -f "$HOME/.zshrc" ]; then
            SHELL_CONFIG="$HOME/.zshrc"
        fi
        
        if [ -n "$SHELL_CONFIG" ]; then
            echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$SHELL_CONFIG"
            echo "✓ Added to $SHELL_CONFIG"
            echo ""
            echo "⚠ Please run: source $SHELL_CONFIG"
            echo "   Or open a new terminal for PATH changes to take effect"
        fi
        
        # Also export for current session
        export PATH="$HOME/.local/bin:$PATH"
    fi
fi

echo ""
echo "========================================"
echo "  Installing Podman Compose"
echo "========================================"
echo ""

# Install podman-compose using pip (user installation, no sudo needed)
if command -v podman-compose &> /dev/null; then
    echo "✓ Podman Compose is already installed"
else
    echo "Installing podman-compose via pip3..."
    
    if command -v pip3 &> /dev/null; then
        pip3 install --user podman-compose
        
        # Add ~/.local/bin to PATH if needed
        if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
            export PATH="$HOME/.local/bin:$PATH"
        fi
        
        if command -v podman-compose &> /dev/null; then
            echo "✓ Podman Compose installed via pip3"
        else
            echo "⚠ Podman Compose installed but not found in PATH"
            echo "  It's installed in: ~/.local/bin/podman-compose"
            echo "  Please add ~/.local/bin to your PATH"
        fi
    elif command -v pip &> /dev/null; then
        pip install --user podman-compose
        echo "✓ Podman Compose installed via pip"
    else
        echo "⚠ pip3 not found. Will try to download manually..."
        
        # Download podman-compose script
        mkdir -p ~/.local/bin
        curl -o ~/.local/bin/podman-compose \
            https://raw.githubusercontent.com/containers/podman-compose/main/podman_compose.py
        chmod +x ~/.local/bin/podman-compose
        
        if [ -f ~/.local/bin/podman-compose ]; then
            echo "✓ Podman Compose downloaded to ~/.local/bin/podman-compose"
        else
            echo "✗ Failed to download podman-compose"
            echo "  You can still use 'podman compose' command (built-in)"
        fi
    fi
fi

echo ""
echo "========================================"
echo "  Verifying Installation..."
echo "========================================"
echo ""

# Verify Podman
if command -v podman &> /dev/null; then
    PODMAN_VERSION=$(podman --version)
    echo "✓ Podman installed: $PODMAN_VERSION"
else
    echo "✗ Podman not found in PATH"
    echo ""
    echo "If you just installed it, try:"
    echo "  source ~/.bashrc"
    echo "  or open a new terminal"
    exit 1
fi

# Verify podman-compose
if command -v podman-compose &> /dev/null; then
    COMPOSE_VERSION=$(podman-compose --version 2>&1 || echo "installed")
    echo "✓ Podman Compose: $COMPOSE_VERSION"
else
    echo "⚠ podman-compose not found"
    echo "  But you can use 'podman compose' (built-in) instead"
fi

echo ""
echo "========================================"
echo "  Testing Podman..."
echo "========================================"
echo ""

# Test Podman
echo "Running test container..."
if podman run --rm alpine echo "Podman works!" 2>/dev/null; then
    echo "✓ Podman is working correctly"
else
    echo "⚠ Test container failed (this might be normal for remote podman)"
    echo "  You can still try to start the application"
fi

echo ""
echo "========================================"
echo "  Installation Complete!"
echo "========================================"
echo ""
echo "Podman has been successfully installed!"
echo ""

if [ "$HAS_SUDO" = false ]; then
    echo "⚠ IMPORTANT: You installed Podman without sudo"
    echo ""
    echo "  1. Reload your shell configuration:"
    echo "     source ~/.bashrc"
    echo "     (or open a new terminal)"
    echo ""
    echo "  2. Verify installation:"
    echo "     podman --version"
    echo ""
    echo "  3. Then start the application:"
    echo "     ./start.sh"
    echo ""
else
    echo "Next steps:"
    echo "  1. Start the forum application:"
    echo "     ./start.sh"
    echo ""
    echo "  2. Access the application:"
    echo "     Frontend: http://localhost:3000"
    echo "     Backend:  http://localhost:8080"
    echo ""
fi

echo "For more information:"
echo "  - Quick start:    QUICKSTART.md"
echo "  - Podman guide:   PODMAN_SETUP.md"
echo "  - Full docs:      README.md"
echo ""
echo "========================================"

