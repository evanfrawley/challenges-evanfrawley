#!/usr/bin/env bash
docker pull evanfrawley/chat-client
docker rm -f chat-client
docker run -d \
--name chat-client \
-p 80:80 -p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
evanfrawley/chat-client
