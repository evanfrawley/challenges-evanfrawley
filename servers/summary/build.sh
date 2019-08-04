#!/usr/bin/env bash

GOOS=linux go build
docker build -t evanfrawley/summarysvc .
docker push evanfrawley/summarysvc
go clean
