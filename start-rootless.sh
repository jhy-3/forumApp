#!/bin/bash

# Start forum application using pure rootless Podman
# No socket, no daemon, no sudo required

set -e

echo "========================================"
echo "  Forum Application - Rootless Start"
echo "========================================"
echo ""

# Set environment for rootless operation
export PODMAN_IGNORE_CGROUPSV1_WARNING=1
export CONTAINER_HOST=""
export XDG_RUNTIME_DIR="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}"

# Ensure runtime directory exists
mkdir -p "$XDG_RUNTIME_DIR" 2>/dev/null || true

echo "Using Podman in rootless mode..."
echo "Runtime directory: $XDG_RUNTIME_DIR"
echo ""

# Check if podman is available
if ! command -v podman &> /dev/null; then
    echo "Error: Podman not found in PATH"
    echo "Please run: source ~/.bashrc"
    exit 1
fi

# Try to get podman info
echo "Testing Podman connection..."
if ! podman info > /dev/null 2>&1; then
    echo "⚠ Podman not fully initialized"
    echo ""
    echo "Attempting to initialize..."
    
    # Try to pull a test image to initialize storage
    podman pull alpine 2>/dev/null || true
fi

echo ""
echo "Starting services with docker-compose format..."
echo ""

# Instead of using podman-compose, use podman directly with compose file
# This bypasses the socket requirement

# Read and parse docker-compose.yml manually
# For MVP, we'll start services one by one

echo "Starting MariaDB..."
podman run -d \
    --name forum-database \
    --network host \
    -e MYSQL_ROOT_PASSWORD=rootpass123 \
    -e MYSQL_DATABASE=forumdb \
    -e MYSQL_USER=forumuser \
    -e MYSQL_PASSWORD=forumpass123 \
    -v forum-mariadb-data:/var/lib/mysql \
    -v "$(pwd)/database/init.sql:/docker-entrypoint-initdb.d/init.sql:ro" \
    mariadb:11.2

echo "Waiting for database to initialize..."
sleep 10

echo ""
echo "Building backend..."
cd backend
podman build -t forum-backend .
cd ..

echo ""
echo "Starting backend..."
podman run -d \
    --name forum-backend \
    --network host \
    -e PORT=8080 \
    -e DB_HOST=127.0.0.1 \
    -e DB_PORT=3306 \
    -e DB_USER=forumuser \
    -e DB_PASSWORD=forumpass123 \
    -e DB_NAME=forumdb \
    -e JWT_SECRET=your-super-secret-jwt-key-change-this-in-production \
    -e CORS_ALLOWED_ORIGINS=http://localhost:3000 \
    forum-backend

echo "Waiting for backend to start..."
sleep 5

echo ""
echo "Building frontend..."
cd frontend
podman build -t forum-frontend .
cd ..

echo ""
echo "Starting frontend..."
podman run -d \
    --name forum-frontend \
    --network host \
    forum-frontend

echo ""
echo "========================================"
echo "  Services Started!"
echo "========================================"
echo ""
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8080"
echo ""
echo "To view logs:"
echo "  podman logs forum-backend"
echo "  podman logs forum-frontend"
echo "  podman logs forum-database"
echo ""
echo "To stop:"
echo "  ./stop-rootless.sh"
echo ""
echo "========================================"

