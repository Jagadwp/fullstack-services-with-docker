#!/bin/bash
set -e

echo "🛑 Stopping Fullstack Assignment"

cd react-frontend || exit
docker compose down
cd ..

cd go-scheduler || exit
# remove -v if you want to keep processed files
docker compose down -v
cd ..

cd python-api || exit
# remove -v if you want to keep received files
docker compose down -v
cd ..

cd php-api || exit
docker compose down
cd ..

echo ""
echo "✅ All services stopped"
