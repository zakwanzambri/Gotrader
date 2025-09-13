#!/bin/bash

# Build script for GoTrader application

echo "Building GoTrader application..."

# Build the backend
echo "Building Go backend..."
go build -o gotrader ./cmd/server/
if [ $? -eq 0 ]; then
    echo "✓ Backend built successfully"
else
    echo "✗ Backend build failed"
    exit 1
fi

# Build the frontend
echo "Building React frontend..."
cd frontend
npm run build
if [ $? -eq 0 ]; then
    echo "✓ Frontend built successfully"
else
    echo "✗ Frontend build failed"
    exit 1
fi

cd ..
echo "✓ GoTrader application built successfully!"
echo ""
echo "To run the application:"
echo "1. Copy .env.example to .env and configure your settings"
echo "2. Run: ./gotrader"
echo "3. Open http://localhost:8080 in your browser"