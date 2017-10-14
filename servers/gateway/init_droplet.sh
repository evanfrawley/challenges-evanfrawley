#!/usr/bin/env bash
# Constants
DOMAIN=evan.gg
SSH_FINGERPRINT="30:12:c7:b9:5a:1b:02:2f:2a:0e:8a:bb:c0:1a:65:e9"
TEMP_IP="192.81.208.97"
# Get region from list
echo "-- Choose your region."
doctl compute region list
read region

# Get Domain
echo "-- What is you domain?"
read domain

# Get Subdomain
echo "-- What is your desired subdomain?"
read subdomain

# Get email
echo "-- What is your email?"
read email

# Get Size
echo "-- What is the desired droplet size?"
doctl compute size list
read size

# Get Name
echo "-- What is the desired droplet name?"
read name
echo "-- Initializing Droplet ..."
echo "-- Please wait ..."

dropletip=$(doctl compute droplet create $name \
--image docker-16-04 \
--region $region \
--size $size \
--ssh-keys $SSH_FINGERPRINT \
--output json \
--wait | jq -r '.[0].networks.v4[0].ip_address')

echo $dropletip

echo "-- Adding A record for this droplet on $domain..."
doctl compute domain records create $domain --record-data $dropletip \
--record-name $subdomain \
--record-type A \
--output json | jq -r '.[0].id' | echo "--Record added with id: $1"

echo "-- SSH into new new droplet ..."
ssh -oStrictHostKeyChecking=no root@$dropletip 'bash -s' < \
init_letsencrypt.sh \
$subdomain \
$domain \
$email
