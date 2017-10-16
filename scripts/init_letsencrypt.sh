#!/usr/bin/env bash
# Args: subdomain, base domain, email

echo "-- Enabling firewall ..."
sudo ufw allow 80
sudo ufw allow 443
sudo ufw allow OpenSSH
"y" | sudo ufw enable

# Install letsencrypt
echo "-- Installing letsencrypt ..."
apt update && apt install -y letsencrypt

echo "-- Obtaining certs ..."
sudo letsencrypt certonly --standalone \
-n \
--agree-tos \
--email $3 \
-d $1.$2
