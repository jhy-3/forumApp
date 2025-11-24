#!/bin/bash

# Forum Application Startup Script
# Supports both Docker and Podman

set -e

echo "======================================"
echo "  Forum Application - Starting..."
echo "======================================"
echo ""

# Determine which container engine to use (Podman preferred)
CONTAINER_ENGINE=""
COMPOSE_CMD=""

if command -v podman &> /dev/null; then
    CONTAINER_ENGINE="podman"
    echo "✓ Found Podman (recommended)"
    
    # Check for podman-compose
    if command -v podman-compose &> /dev/null; then
        COMPOSE_CMD="podman-compose"
    else
        echo "Warning: podman-compose not found. Please install it:"
        echo "  sudo apt-get install podman-compose"
        echo "  or: pip3 install podman-compose"
        exit 1
    fi
elif command -v docker &> /dev/null; then
    CONTAINER_ENGINE="docker"
    echo "✓ Found Docker"
    
    # Check for docker compose
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        echo "Error: Docker Compose is not installed."
        exit 1
    fi
else
    echo "Error: Neither Podman nor Docker is installed."
    echo ""
    echo "Recommended: Install Podman (more secure, no root required)"
    echo "  Ubuntu/Debian: sudo apt-get install podman podman-compose"
    echo "  Fedora/RHEL:   sudo dnf install podman podman-compose"
    echo ""
    echo "See PODMAN_SETUP.md for detailed instructions"
    exit 1
fi

echo "Using: $COMPOSE_CMD"
echo "Starting services..."
$COMPOSE_CMD up -d

echo ""
echo "Waiting for services to be ready..."
sleep 5

# Check if backend is ready
echo -n "Checking backend health..."
for i in {1..30}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo " ✓ Ready"
        break
    fi
    if [ $i -eq 30 ]; then
        echo " ✗ Timeout"
        echo "Backend is not responding. Check logs with: docker-compose logs backend"
        exit 1
    fi
    sleep 2
    echo -n "."
done

echo ""
echo "======================================"
echo "  Services are ready!"
echo "======================================"
echo ""
echo "  Frontend:  http://localhost:3000"
echo "  Backend:   http://localhost:8080"
echo "  Database:  localhost:3306"
echo ""
echo "Default forums have been created:"
echo "  - General Discussion"
echo "  - Announcements"
echo "  - Technical Support"
echo "  - Off Topic"
echo ""
echo "To view logs:       docker-compose logs -f"
echo "To stop services:   docker-compose down"
echo "To stop & clean:    docker-compose down -v"
echo ""
echo "======================================"

