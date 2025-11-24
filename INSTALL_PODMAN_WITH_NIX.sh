#!/bin/bash

# Install Podman using Nix package manager (no sudo required)
# Nix can install a complete Podman in user space

set -e

echo "========================================"
echo "  Installing Podman via Nix"
echo "  (No sudo required!)"
echo "========================================"
echo ""

# Step 1: Check if Nix is installed
if command -v nix &> /dev/null; then
    echo "✓ Nix is already installed"
else
    echo "[1/3] Installing Nix package manager..."
    echo ""
    echo "Nix is a powerful package manager that works without root."
    echo "It will install to: ~/.nix-profile/"
    echo ""
    
    # Install Nix in single-user mode (no sudo needed)
    curl -L https://nixos.org/nix/install | sh -s -- --no-daemon
    
    # Source nix
    if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then
        . ~/.nix-profile/etc/profile.d/nix.sh
    fi
    
    echo ""
    echo "✓ Nix installed successfully!"
fi

# Ensure Nix is in PATH
if [ -e ~/.nix-profile/etc/profile.d/nix.sh ]; then
    . ~/.nix-profile/etc/profile.d/nix.sh
fi

echo ""
echo "[2/3] Installing Podman via Nix..."
echo "This may take a few minutes..."
echo ""

# Install Podman using Nix
nix-env -iA nixpkgs.podman

echo ""
echo "[3/3] Installing podman-compose..."

# Install podman-compose
if command -v pip3 &> /dev/null; then
    pip3 install --user podman-compose
else
    # Install via Nix
    nix-env -iA nixpkgs.podman-compose
fi

echo ""
echo "========================================"
echo "  Verifying Installation"
echo "========================================"
echo ""

# Verify Podman
if command -v podman &> /dev/null; then
    PODMAN_VERSION=$(podman --version)
    echo "✓ Podman: $PODMAN_VERSION"
else
    echo "✗ Podman not found"
    echo "  Try: source ~/.nix-profile/etc/profile.d/nix.sh"
    exit 1
fi

# Verify podman-compose
if command -v podman-compose &> /dev/null; then
    echo "✓ Podman Compose: installed"
else
    echo "⚠ Podman Compose not found (you can use 'podman compose' instead)"
fi

echo ""
echo "========================================"
echo "  Testing Podman"
echo "========================================"
echo ""

# Test Podman
echo "Running test container..."
if podman run --rm alpine echo "Podman works!"; then
    echo "✓ Podman is working correctly!"
else
    echo "⚠ Test failed, but you can still try to use it"
fi

echo ""
echo "========================================"
echo "  Installation Complete!"
echo "========================================"
echo ""
echo "Podman has been installed via Nix!"
echo ""
echo "Next steps:"
echo "  1. Source Nix profile:"
echo "     source ~/.nix-profile/etc/profile.d/nix.sh"
echo "     (or add it to your ~/.bashrc)"
echo ""
echo "  2. Start the application:"
echo "     ./start.sh"
echo ""
echo "  3. Access at:"
echo "     Frontend: http://localhost:3000"
echo "     Backend:  http://localhost:8080"
echo ""
echo "Note: Nix installs are stored in /nix/store/"
echo "      User profile links in ~/.nix-profile/"
echo ""
echo "========================================"

