#!/usr/bin/env bash
subdomain=
if [ $1 ] ; then
    subdomain=$1
else
    subdomain="api"
fi
export TLSCERT=/etc/letsencrypt/live/${subdomain}.evan.gg/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/${subdomain}.evan.gg/privkey.pem

# api-network addresses
export MONGO_SERVER_NAME=mongosvr:27017
export REDIS_SERVER_NAME=redissvr:6379
export MESSAGING_SERVICE_NAME=messagingsvc
export SUMMARY_SERVICE_NAME=summarysvc
export DEFAULT_NETWORK_PORT=":80"
export GATEWAY_API_NAME=gateway-api

export LOCALHOST=localhost

export DOCKER_NETWORK=api-network

# Mongo DB
# Init Mongo store
docker rm -f $MONGO_SERVER_NAME
docker run -d \
--name $MONGO_SERVER_NAME \
--network $DOCKER_NETWORK \
mongo

# Redis Store
# Init Redis store
docker rm -f $REDIS_SERVER_NAME
docker run -d \
--name $REDIS_SERVER_NAME \
--network $DOCKER_NETWORK \
redis

# Set up messaging microservice
docker rm -f messagingsvc
docker pull evanfrawley/messagingsvc
docker run -d \
--name $MESSAGING_SERVICE_NAME \
--network $DOCKER_NETWORK \
-e DBADDR=$MONGO_SERVER_NAME \
-e NODE_ADDR=$LOCALHOST$DEFAULT_NETWORK_PORT \
evanfrawley/messagingsvc

# Set up summary microservice
docker rm -f summarysvc
docker pull evanfrawley/summarysvc
docker run -d \
--name $SUMMARY_SERVICE_NAME \
--network $DOCKER_NETWORK \
-e SUMMARY_ADDR=$LOCALHOST$DEFAULT_NETWORK_PORT \
evanfrawley/summarysvc

# Set up gateway
docker rm -f gateway-api
docker pull evanfrawley/gateway-api
docker run -d \
--name $GATEWAY_API_NAME \
--network $DOCKER_NETWORK \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e GOADDR=gateway-api:443 \
-e MONGO_ADDR=$MONGO_SERVER_NAME \
-e REDIS_ADDR=$REDIS_SERVER_NAME \
-e MSGSVC_ADDRS=$MESSAGING_SERVICE_NAME$DEFAULT_NETWORK_PORT \
-e SUMMARYSVC_ADDRS=$SUMMARY_SERVICE_NAME$DEFAULT_NETWORK_PORT \
evanfrawley/gateway-api
