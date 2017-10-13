#!/usr/bin/env bash
docker push evanfrawley/chat-client
ssh -oStrictHostKeyChecking=no root@159.203.116.26 'bash -s' < run.sh
