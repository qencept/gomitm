#!/bin/bash

if [ ! -f ca_mitm.crt ] || [ ! -f ca_mitm.key ]; then
    ./generate_scripts.sh
fi

sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner root --dport 80 -j REDIRECT --to-port 8888
sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner root --dport 443 -j REDIRECT --to-port 8888

go build -o gomitm ../../cmd/gomitm
sudo ./gomitm --config config.yaml
