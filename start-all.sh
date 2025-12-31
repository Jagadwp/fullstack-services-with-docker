#!/bin/bash
set -e

echo "🚀 Starting Fullstack Assignment"

./setup-network.sh

echo ""
echo "▶️ Starting PHP API"
cd php-api
docker compose up -d --build
cd ..

echo ""
echo "▶️ Starting Python API"
cd python-api
docker compose up -d --build
cd ..

echo ""
echo "▶️ Starting Go Scheduler"
cd go-scheduler
docker compose up -d --build
cd ..

echo ""
echo "▶️ Starting React Frontend"
cd react-frontend
docker compose build --no-cache && docker compose up -d
cd ..

echo ""
echo "✅ All services are up and running"
