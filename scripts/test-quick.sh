#!/bin/bash
# Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
# Quick test runner for development (no coverage)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🚀 Running quick tests (no coverage)..."
echo ""

# Backend
echo "📦 Backend tests..."
cd "$PROJECT_ROOT/backend"
go test -short ./... -timeout=30s

echo ""
echo "✅ Quick tests passed!"
