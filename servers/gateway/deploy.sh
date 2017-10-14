#!/usr/bin/env bash
source build.sh
docker push evanfrawley/gateway-api

ssh -oStrictHostKeyChecking=no root@198.211.98.217 'bash -s' < run.sh
