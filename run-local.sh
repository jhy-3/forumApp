#!/bin/bash

# Run Forum Application Locally (No Container Required)

echo "========================================"
echo "  Forum Application - Local Mode"
echo "  (No Docker/Podman needed!)"
echo "========================================"
echo ""

# Check Go
echo "[1/4] Checking Go..."
if command -v go &> /dev/null; then
    echo "  ✓ Go installed: $(go version)"
else
    echo "  ✗ Go not found"
    echo ""
    echo "  Install Go without sudo:"
    echo "    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz"
    echo "    tar -xzf go1.21.0.linux-amd64.tar.gz -C ~/"
    echo "    echo 'export PATH=\$HOME/go/bin:\$PATH' >> ~/.bashrc"
    echo "    source ~/.bashrc"
    exit 1
fi

# Check Node.js
echo ""
echo "[2/4] Checking Node.js..."
if command -v node &> /dev/null; then
    echo "  ✓ Node.js installed: $(node --version)"
else
    echo "  ✗ Node.js not found"
    echo ""
    echo "  Install Node.js without sudo:"
    echo "    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash"
    echo "    source ~/.bashrc"
    echo "    nvm install 18"
    exit 1
fi

# Check database configuration
echo ""
echo "[3/4] Checking database configuration..."
if [ -f backend/.env ]; then
    echo "  ✓ Backend .env found"
    DB_HOST=$(grep DB_HOST backend/.env | cut -d= -f2)
    echo "  Database host: $DB_HOST"
else
    echo "  ⚠ No .env file found, using default"
    cp backend/.env.example backend/.env 2>/dev/null || true
fi

echo ""
echo "========================================"
echo "  Database Setup Required"
echo "========================================"
echo ""
echo "Before starting, you need a MySQL/MariaDB database."
echo ""
echo "Options:"
echo "  1. Use Railway.app (free, no sudo) - Recommended"
echo "     https://railway.app"
echo ""
echo "  2. Use local MySQL (if installed)"
echo "     mysql -u root -p < database/init.sql"
echo ""
echo "  3. Use PlanetScale (free)"
echo "     https://planetscale.com"
echo ""
read -p "Have you configured the database? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "Please configure database first."
    echo "See: 本地运行指南.md"
    exit 1
fi

echo ""
echo "[4/4] Starting services..."
echo ""

# Create log directory
mkdir -p logs

# Start backend in background
echo "Starting backend..."
cd backend
go mod download 2>/dev/null || true
nohup go run main.go > ../logs/backend.log 2>&1 &
BACKEND_PID=$!
echo "  Backend PID: $BACKEND_PID"
cd ..

# Wait for backend to start
echo "  Waiting for backend to start..."
sleep 3

# Test backend
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "  ✓ Backend is running"
else
    echo "  ⚠ Backend may not be ready yet"
    echo "  Check logs: tail -f logs/backend.log"
fi

# Start frontend
echo ""
echo "Starting frontend..."
cd frontend

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "  Installing frontend dependencies..."
    npm install
fi

echo "  Starting React dev server..."
echo "  (This will open your browser automatically)"
echo ""

# Frontend runs in foreground
npm start

# Cleanup on exit
echo ""
echo "Shutting down backend..."
kill $BACKEND_PID 2>/dev/null || true
echo "Done!"

