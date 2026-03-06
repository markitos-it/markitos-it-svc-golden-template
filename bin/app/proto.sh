#!/usr/bin/env sh
set -eu

echo "== protoc version =="
protoc --version

echo "== protoc-gen-go =="
command -v protoc-gen-go >/dev/null 2>&1

echo "== protoc-gen-go-grpc =="
command -v protoc-gen-go-grpc >/dev/null 2>&1

echo "== generating protobuf go files =="
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/golden.proto

echo "== done =="
ls -la proto