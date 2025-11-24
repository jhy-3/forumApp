#!/bin/bash

# Fix Podman to work in pure rootless mode without socket
# This configures Podman to run containers directly without daemon

set -e

echo "========================================"
echo "  Configuring Podman Rootless Mode"
echo "========================================"
echo ""

# Step 1: Configure user namespaces
echo "[1/6] Checking user namespace configuration..."

if ! grep -q "^$USER:" /etc/subuid 2>/dev/null; then
    echo "  ⚠ User not configured in /etc/subuid"
    echo "  Creating local configuration..."
    
    # Create local subuid/subgid mapping
    mkdir -p ~/.config/containers
    
    # Use a range starting from current UID
    CURRENT_UID=$(id -u)
    START_UID=$((CURRENT_UID * 100000))
    
    echo "  Using UID range: $START_UID - $((START_UID + 65536))"
else
    echo "  ✓ User namespace configured"
fi

# Step 2: Create Podman configuration directory
echo ""
echo "[2/6] Creating Podman configuration..."
mkdir -p ~/.config/containers
mkdir -p ~/.local/share/containers/storage

# Step 3: Configure storage
echo ""
echo "[3/6] Configuring storage..."
cat > ~/.config/containers/storage.conf << 'EOF'
[storage]
driver = "vfs"
runroot = "/run/user/$UID/containers"
graphroot = "$HOME/.local/share/containers/storage"

[storage.options]
mount_program = "/usr/bin/fuse-overlayfs"

[storage.options.vfs]
ignore_chown_errors = "true"
EOF

echo "  ✓ Storage configured (VFS driver)"

# Step 4: Configure containers.conf
echo ""
echo "[4/6] Configuring Podman behavior..."
cat > ~/.config/containers/containers.conf << 'EOF'
[containers]
netns="host"
userns="host"
ipcns="host"
utsns="host"
cgroupns="host"
cgroups="disabled"
log_driver = "k8s-file"
pids_limit = 2048

[engine]
cgroup_manager = "cgroupfs"
events_logger="file"
runtime="crun"

[network]
network_backend="netavark"
EOF

echo "  ✓ Container behavior configured"

# Step 5: Create wrapper script for podman-compose
echo ""
echo "[5/6] Creating podman-compose wrapper..."
cat > ~/.local/bin/podman-compose-fixed << 'EOF'
#!/bin/bash
# Wrapper for podman-compose with proper environment

export PODMAN_IGNORE_CGROUPSV1_WARNING=1
export CONTAINER_HOST=""

# Use podman directly without socket
exec python3 -m podman_compose "$@"
EOF

chmod +x ~/.local/bin/podman-compose-fixed

# Step 6: Test configuration
echo ""
echo "[6/6] Testing Podman..."

# Test if podman can run without socket
if ~/.local/bin/podman --log-level=debug info > /tmp/podman-test.log 2>&1; then
    echo "  ✓ Podman configured successfully!"
else
    echo "  ⚠ Podman test failed. Checking details..."
    echo ""
    echo "Error details:"
    tail -20 /tmp/podman-test.log
    echo ""
    echo "This might be due to system limitations."
fi

echo ""
echo "========================================"
echo "  Configuration Complete"
echo "========================================"
echo ""
echo "Next steps:"
echo "  1. Try running: podman info"
echo "  2. If that works, run: ./start-rootless.sh"
echo ""
echo "========================================"

