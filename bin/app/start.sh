#!/bin/bash

set -e

cd "$(dirname "$0")/../.."

make proto

export GRPC_PORT=${GRPC_PORT:-3000}
export GRPC_TLS_ENABLED=${GRPC_TLS_ENABLED:-false}
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-55432}
export DB_USER=${DB_USER:-markitos-it-svc-goldens}
export DB_PASS=${DB_PASS:-markitos-it-svc-goldens}
export DB_NAME=${DB_NAME:-markitos-it-svc-goldens}

echo "🚀 Starting markitos-it-svc-goldens (Go)..."
echo "📡 GRPC_PORT.......: $GRPC_PORT"
echo "📦 DB_HOST.........: $DB_HOST:$DB_PORT"
echo "📦 DB_USER.........: $DB_USER"
echo "📦 DB_NAME.........: $DB_NAME"
echo "📡 GRPC_TLS_ENABLED: $GRPC_TLS_ENABLED"
echo ""

docker compose up