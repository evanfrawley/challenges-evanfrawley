#!/usr/bin/env bash
subdomain=
if [ $1 ] ; then
    subdomain=$1
else
    subdomain="api"
fi
export TLSCERT=/etc/letsencrypt/live/${subdomain}.evan.gg/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/${subdomain}.evan.gg/privkey.pem

docker rm -f gateway-api
docker pull evanfrawley/gateway-api
docker run -d \
--name gateway-api \
--network api-network \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e GO_ADDR=gateway-api:443 \
evanfrawley/gateway-api
