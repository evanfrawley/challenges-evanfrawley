#!/usr/bin/env bash
set -e
dropletName=
if [ $1 ] ; then
    dropletName=$1
else
    dropletName="chat-client"
fi

dropletIp=$(
    doctl compute droplet list \
    --output json \
    | jq -r --arg n ${dropletName} \
        '.[] as $i | $i.name | scan($n) | $i.networks.v4[0].ip_address'
)


ssh -oStrictHostKeyChecking=no root@${dropletIp} 'bash -s' < run.sh
