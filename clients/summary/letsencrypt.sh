#!/usr/bin/env bash
sudo ufw allow 80
sudo ufw allow 443
sudo ufw allow OpenSSH
ufw enable

sudo letsencrypt certonly --standalone \
-n \
--agree-tos \
--email frawley@uw.edu \
-d chat-client.evan.gg
