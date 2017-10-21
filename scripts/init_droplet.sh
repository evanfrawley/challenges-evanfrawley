#!/usr/bin/env bash
# This script will:
# - initialize a Droplet in DigitalOcean with custom configuration
# - set up a subdomain A record on DigitalOcean
# - enable letsencrypt / HTTPS on the new droplet

# Constants
SSH_FINGERPRINT="30:12:c7:b9:5a:1b:02:2f:2a:0e:8a:bb:c0:1a:65:e9"

# Get Domain
echo "-- What domain would you like to add to?"
read domain

# Get Name
echo "-- What is the desired droplet name?"
read name

existingDropletName=$(
    doctl compute droplet list \
        --output json \
        | jq -r --arg n ${name} '.[].name | scan($n)'
)

dropletNameExists=
case ${existingDropletName} in
    ${name}) dropletNameExists=true ;;
    "") dropletNameExists=false ;;
    *) echo ERROR: failed to parse doctl output ; exit 1 ;;
esac

dropletIp=
if [ ${existingDropletName} ] ; then
    dropletIp=$(
        doctl compute droplet list \
            --output json \
            | jq -r --arg n ${name} \
                '.[] as $i | $i.name | scan($n) | $i.networks.v4[0].ip_address'
    )
    echo "Seems like a droplet with this name is already set up on ip: ${dropletIp}"
else
    # Get email
    echo "-- What is your email?"
    read email

    # Get region from list
    echo "-- Choose your region."
    doctl compute region list
    read region

    # Get Size
    echo "-- What is the desired droplet size?"
    doctl compute size list
    read size

    echo "-- Initializing Droplet ..."
    echo "-- Please wait ..."

    dropletIp=$(
        doctl compute droplet create ${name} \
            --image docker-16-04 \
            --region ${region} \
            --size ${size} \
            --ssh-keys ${SSH_FINGERPRINT} \
            --output json \
            --wait \
            | jq -r \
                '.[0].networks.v4[0].ip_address'
    )
    # Get Subdomain
    echo "-- What is your desired subdomain?"
    read subdomain

    recordName=$(
        doctl compute domain records list ${domain} \
            --output json \
            | jq -r --arg sd ${subdomain} \
                '.[].name | scan($sd)'
    )

    subdomainRecordExists=
    case ${recordName} in
        ${subdomain}) subdomainRecordExists=true ;;
        "") subdomainRecordExists=false ;;
        *) echo ERROR: failed to parse doctl output 1; exit 1 ;;
    esac

    if [ ${subdomainRecordExists} == "false" ] ; then
        echo "-- Adding A record for this droplet on $domain..."
        doctl compute domain records create ${domain} --record-data ${dropletIp} \
            --record-name ${subdomain} \
            --record-type A \
            --output json \
            | jq -r '.[0].id'
    fi

    echo "-- SSH into new new droplet ..."
    ssh -oStrictHostKeyChecking=no root@${dropletIp} 'bash -s' < \
        init_letsencrypt.sh \
            ${subdomain} \
            ${domain} \
            ${email}
fi
