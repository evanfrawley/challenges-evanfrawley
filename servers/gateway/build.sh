#!/usr/bin/env bash
GOOS=linux go build
docker build -t evanfrawley/gateway-api .
go clean
