#!/usr/bin/env bash

# Init Mongo store
docker rm -f mongosvr
docker run -d \
--name mongosvr \
--network api-network \
mongo
