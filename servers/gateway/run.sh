#!/usr/bin/env bash
export TLSCERT=/etc/letsencrypt/live/test2.evan.gg/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/test2.evan.gg/privkey.pem

docker rm -f 344gateway

docker run -d \
--name 344gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
evanfrawley/gateway-api
