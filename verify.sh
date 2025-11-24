#!/bin/bash

# Verification script to check if everything is ready

echo "========================================"
echo "  Forum Application - Environment Check"
echo "========================================"
echo ""

# Check 1: Podman binary
echo "[1/5] Checking Podman binary..."
if [ -f ~/.local/bin/podman ]; then
    echo "  ✓ Podman binary found at ~/.local/bin/podman"
    ls -lh ~/.local/bin/podman
else
    echo "  ✗ Podman binary not found"
    echo "  Please run: ./INSTALL_PODMAN.sh"
fi

echo ""

# Check 2: Podman Compose
echo "[2/5] Checking Podman Compose..."
if [ -f ~/.local/bin/podman-compose ]; then
    echo "  ✓ Podman Compose found"
else
    echo "  ✗ Podman Compose not found"
fi

echo ""

# Check 3: PATH configuration
echo "[3/5] Checking PATH configuration..."
if grep -q ".local/bin" ~/.bashrc; then
    echo "  ✓ PATH configured in ~/.bashrc"
    grep "local/bin" ~/.bashrc | tail -1
else
    echo "  ✗ PATH not configured"
    echo "  Please run: ./INSTALL_PODMAN.sh"
fi

echo ""

# Check 4: Current PATH
echo "[4/5] Checking current PATH..."
if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
    echo "  ✓ ~/.local/bin is in current PATH"
else
    echo "  ⚠ ~/.local/bin NOT in current PATH"
    echo "  Please run: source ~/.bashrc"
fi

echo ""

# Check 5: Podman availability
echo "[5/5] Checking Podman availability..."
if command -v podman &> /dev/null; then
    echo "  ✓ Podman is available!"
    podman --version
else
    echo "  ⚠ Podman command not found"
    echo ""
    echo "  To fix this, run:"
    echo "    source ~/.bashrc"
    echo "  Or open a new terminal window"
fi

echo ""
echo "========================================"
echo "  Summary"
echo "========================================"
echo ""

# Overall status
ALL_OK=true

if [ ! -f ~/.local/bin/podman ]; then
    ALL_OK=false
fi

if ! grep -q ".local/bin" ~/.bashrc 2>/dev/null; then
    ALL_OK=false
fi

if $ALL_OK; then
    if command -v podman &> /dev/null; then
        echo "✓ Everything is ready!"
        echo ""
        echo "  You can now start the application:"
        echo "    ./start.sh"
        echo ""
    else
        echo "⚠ Almost ready!"
        echo ""
        echo "  Next step: Reload your shell configuration"
        echo "    source ~/.bashrc"
        echo ""
        echo "  Then start the application:"
        echo "    ./start.sh"
        echo ""
    fi
else
    echo "✗ Setup incomplete"
    echo ""
    echo "  Please run the installation script:"
    echo "    ./INSTALL_PODMAN.sh"
    echo ""
fi

echo "========================================"

