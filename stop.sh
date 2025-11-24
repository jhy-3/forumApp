#!/bin/bash

# Forum Application Stop Script
# Supports both Docker and Podman

echo "======================================"
echo "  Forum Application - Stopping..."
echo "======================================"
echo ""

# Determine which compose command to use
COMPOSE_CMD=""

if command -v podman-compose &> /dev/null; then
    COMPOSE_CMD="podman-compose"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
else
    echo "Error: No compose command found."
    exit 1
fi

echo "Using: $COMPOSE_CMD"
$COMPOSE_CMD down

echo ""
echo "Services stopped successfully!"
echo ""
echo "To remove all data: $COMPOSE_CMD down -v"
echo "======================================"

