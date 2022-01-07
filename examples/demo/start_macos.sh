#!/bin/bash

if [ ! -f ca_mitm.crt ] || [ ! -f ca_mitm.key ]; then
  ./generate_certs.sh
fi

if [ ! -f doh ] || [ ! -f http ] || [ ! -f session ]; then
  mkdir -p doh
  mkdir -p http
  mkdir -p session
fi

sudo sysctl -w net.inet.ip.forwarding=1
sudo pfctl -ef pf.conf

ulimit -n 65536
sudo go run ../../cmd/gomitm --config config.yaml
