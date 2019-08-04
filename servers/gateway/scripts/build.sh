#!/usr/bin/env bash
cd ~/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway
GOOS=linux go build
docker build -t evanfrawley/gateway-api .
go clean
