#!/usr/bin/env bash

subdomain=
if [ $1 ] ; then
    subdomain=$1
else
    subdomain="api"
fi

dropletIp=$(
    doctl compute droplet list \
    --output json \
    | jq -r --arg n ${subdomain} \
        '.[] as $i | $i.name | scan($n) | $i.networks.v4[0].ip_address'
)

source build.sh
docker push evanfrawley/gateway-api

ssh -oStrictHostKeyChecking=no root@${dropletIp} 'bash -s' < run.sh ${subdomain}
