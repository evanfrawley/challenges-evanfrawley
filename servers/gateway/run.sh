#!/usr/bin/env bash
subdomain=
if [ $1 ] ; then
    subdomain=$1
else
    subdomain="api"
fi
export TLSCERT=/etc/letsencrypt/live/${subdomain}.evan.gg/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/${subdomain}.evan.gg/privkey.pem

docker rm -f 344gateway

docker run -d \
--name 344gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
evanfrawley/gateway-api
