#!/bin/bash

if [ ! -f ca_mitm.crt ] || [ ! -f ca_mitm.key ]; then
    ./generate_certs.sh
fi

go build -o gomitm ../../cmd/gomitm

sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner root --dport 80 -j REDIRECT --to-port 8888
sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner root --dport 443 -j REDIRECT --to-port 8888

sudo ./gomitm --config config.yaml
