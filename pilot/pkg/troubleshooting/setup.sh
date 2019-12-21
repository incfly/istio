#!/bin/bash

protoc --go_out=plugins=grpc:. api/service.proto

func server() {
  go run ./cmd/istiod
}

func client() {
  go run ./cmd/agent
}