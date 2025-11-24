#!/bin/bash

# Start Podman socket service in user mode (no sudo required)
# This allows podman-remote to connect

set -e

echo "========================================"
echo "  Starting Podman User Service"
echo "========================================"
echo ""

# Create systemd user directory
mkdir -p ~/.config/systemd/user

# Create podman.socket service
echo "[1/4] Creating Podman socket service..."
cat > ~/.config/systemd/user/podman.socket << 'EOF'
[Unit]
Description=Podman API Socket
Documentation=man:podman-system-service(1)

[Socket]
ListenStream=%t/podman/podman.sock
SocketMode=0660

[Install]
WantedBy=sockets.target
EOF

# Create podman.service
echo "[2/4] Creating Podman service..."
cat > ~/.config/systemd/user/podman.service << 'EOF'
[Unit]
Description=Podman API Service
Requires=podman.socket
After=podman.socket
Documentation=man:podman-system-service(1)
StartLimitIntervalSec=0

[Service]
Type=exec
KillMode=process
Environment=LOGGING="--log-level=info"
ExecStart=/home/%u/.local/bin/podman $LOGGING system service

[Install]
WantedBy=default.target
EOF

# Replace %u with actual username
sed -i "s|/home/%u|$HOME|g" ~/.config/systemd/user/podman.service

# Reload systemd user daemon
echo "[3/4] Reloading systemd user daemon..."
systemctl --user daemon-reload

# Enable and start the socket
echo "[4/4] Starting Podman socket..."
systemctl --user enable podman.socket 2>/dev/null || true
systemctl --user start podman.socket

echo ""
echo "Waiting for socket to be ready..."
sleep 2

# Check status
if systemctl --user is-active podman.socket >/dev/null 2>&1; then
    echo "✓ Podman socket is running!"
    echo ""
    echo "Socket location: $XDG_RUNTIME_DIR/podman/podman.sock"
    
    # Test the socket
    if podman info >/dev/null 2>&1; then
        echo "✓ Podman is responding!"
    else
        echo "⚠ Socket started but Podman not responding"
    fi
else
    echo "✗ Failed to start Podman socket"
    echo ""
    echo "Checking logs..."
    journalctl --user -u podman.socket -n 20
fi

echo ""
echo "========================================"
echo "  Service Status"
echo "========================================"
echo ""
systemctl --user status podman.socket --no-pager || true

echo ""
echo "========================================"
echo "  Next Steps"
echo "========================================"
echo ""
echo "If the socket is running, you can now use:"
echo "  ./start.sh"
echo ""
echo "To check service status:"
echo "  systemctl --user status podman.socket"
echo ""
echo "To stop the service:"
echo "  systemctl --user stop podman.socket"
echo ""
echo "========================================"

