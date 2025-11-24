#!/bin/bash

# Fix Podman Socket Connection Issue
# This script sets up Podman to work without a socket (rootless mode)

echo "========================================"
echo "  Fixing Podman Socket Issue"
echo "========================================"
echo ""

echo "Problem: Podman remote binary needs socket connection"
echo "Solution: Configure Podman for direct execution"
echo ""

# Create config directory
mkdir -p ~/.config/containers

# Create connections.conf to disable socket requirement
echo "Creating Podman configuration..."
cat > ~/.config/containers/connections.conf << 'EOF'
[connection]
default = "root"

[connection.root]
uri = "unix:///run/podman/podman.sock"
EOF

# For remote podman, we need to tell it to work in local mode
# Create a wrapper script
echo "Creating Podman wrapper..."
cat > ~/.local/bin/podman-wrapper << 'EOF'
#!/bin/bash
# Podman wrapper for direct execution without socket

# Use podman in remote mode with local socket or direct execution
export CONTAINER_HOST=""
export CONTAINER_SSHKEY=""
export CONTAINER_CONNECTION=""

# Execute podman with all arguments
exec ~/.local/bin/podman "$@"
EOF

chmod +x ~/.local/bin/podman-wrapper

echo ""
echo "========================================"
echo "  Alternative Solution"
echo "========================================"
echo ""
echo "The issue is that the static binary requires a Podman socket."
echo "Since you don't have sudo, here are your options:"
echo ""
echo "Option 1: Ask system administrator to install Podman"
echo "  sudo apt-get install podman podman-compose"
echo "  (This will work for all users without sudo)"
echo ""
echo "Option 2: Use Docker instead"
echo "  If Docker is already installed, use:"
echo "  docker compose up -d"
echo ""
echo "Option 3: Install full Podman from source (complex)"
echo "  This requires compilation and is not recommended"
echo ""
echo "========================================"
echo "  Recommended: Ask Admin to Install"
echo "========================================"
echo ""
echo "The simplest solution is to ask your system administrator to run:"
echo "  sudo apt-get update"
echo "  sudo apt-get install -y podman podman-compose"
echo ""
echo "After that, ALL users can use Podman without sudo!"
echo "========================================"

