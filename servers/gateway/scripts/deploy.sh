#!/usr/bin/env bash
set -e
dropletName=
if [ $1 ] ; then
    dropletName=$1
else
    dropletName="api"
fi

dropletIp=$(
    doctl compute droplet list \
    --output json \
    | jq -r --arg n ${dropletName} \
        '.[] as $i | $i.name | scan($n) | $i.networks.v4[0].ip_address'
)

source /Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/scripts/build.sh
docker push evanfrawley/gateway-api

ssh -oStrictHostKeyChecking=no root@${dropletIp} 'bash -s' < /Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/scripts/run.sh ${dropletName}
