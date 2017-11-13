#!/usr/bin/env bash
# Init Redis store
docker rm -f redissvr
docker run -d \
--name redissvr \
--network api-network \
redis
