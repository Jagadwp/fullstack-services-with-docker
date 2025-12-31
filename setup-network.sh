#!/bin/bash
set -e

NETWORK_NAME="fullstack_network"

echo "🔍 Checking Docker network: $NETWORK_NAME"

if docker network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
  echo "✅ Network '$NETWORK_NAME' already exists"
else
  echo "➕ Creating Docker network '$NETWORK_NAME'"
  docker network create "$NETWORK_NAME"
  echo "✅ Network '$NETWORK_NAME' created"
fi

echo ""
docker network inspect "$NETWORK_NAME" | grep Name
