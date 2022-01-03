#!/bin/bash

if [ ! -f ca_mitm.crt ] || [ ! -f ca_mitm.key ]; then
    ./generate_certs.sh
fi

if [ ! -f doh ] || [ ! -f http ] || [ ! -f session ]; then
  mkdir -p doh
  mkdir -p http
  mkdir -p session
fi

sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner root --dport 80 -j REDIRECT --to-port 8888
sudo iptables -t nat -A OUTPUT -p tcp -m owner ! --uid-owner root --dport 443 -j REDIRECT --to-port 8888

sudo go run ../../cmd/gomitm --config config.yaml
