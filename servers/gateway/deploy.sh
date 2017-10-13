#!/usr/bin/env bash
source build.sh
docker push evanfrawley/gateway-api

ssh -oStrictHostKeyChecking=no root@104.236.5.41 'bash -s' < run.sh
