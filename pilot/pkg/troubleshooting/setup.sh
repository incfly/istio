#!/bin/bash

proto() {
  protoc --go_out=plugins=grpc:. api/service.proto
}

server() {
  go run ./cmd/istiod "$@"
}

agent() {
  go run ./cmd/agent "$@"
}

cli() {
  go run ./cmd/cli "$@"
}