#!/bin/bash

set -e

cd "$(dirname "$0")/../.."

echo "🛑 Stopping markitos-it-svc-goldens PostgreSQL..."
docker compose down -v --remove-orphans markitos-it-svc-goldens-postgres markitos-it-svc-goldens > /dev/null 2>&1 || true
echo "✅ markitos-it-svc-goldens PostgreSQL stopped."
echo
