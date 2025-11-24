#!/bin/bash

# Stop all forum containers

echo "========================================"
echo "  Stopping Forum Application"
echo "========================================"
echo ""

# Stop and remove containers
echo "Stopping containers..."
podman stop forum-frontend forum-backend forum-database 2>/dev/null || true

echo "Removing containers..."
podman rm forum-frontend forum-backend forum-database 2>/dev/null || true

echo ""
echo "✓ All containers stopped"
echo ""
echo "To remove data volumes:"
echo "  podman volume rm forum-mariadb-data"
echo ""
echo "========================================"

