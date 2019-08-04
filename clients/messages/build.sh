#!/usr/bin/env bash
yarn build
docker build -t evanfrawley/chat-client .
docker push evanfrawley/chat-client