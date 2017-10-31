#!/usr/bin/env bash
cd /Users/evanfrawley/go/src/github.com/info344-a17/challenges-evanfrawley/servers/gateway/scripts
source init_droplet.sh
cd ../servers/gateway/
source deploy.sh "api"