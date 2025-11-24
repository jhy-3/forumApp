#!/bin/bash

# Stop locally running services

echo "========================================"
echo "  Stopping Local Services"
echo "========================================"
echo ""

# Find and kill Go backend
echo "Stopping backend..."
pkill -f "go run main.go" 2>/dev/null && echo "  ✓ Backend stopped" || echo "  - Backend not running"

# Find and kill Node frontend
echo "Stopping frontend..."
pkill -f "react-scripts start" 2>/dev/null && echo "  ✓ Frontend stopped" || echo "  - Frontend not running"

echo ""
echo "✓ All local services stopped"
echo "========================================"

