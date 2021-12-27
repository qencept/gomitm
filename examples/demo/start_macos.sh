#!/bin/bash

if [ ! -f ca_mitm.crt ] || [ ! -f ca_mitm.key ]; then
    ./generate_certs.sh
fi

go build -o gomitm ../../cmd/gomitm

sudo sysctl -w net.inet.ip.forwarding=1
sudo pfctl -ef pf.conf

sudo ./gomitm --config config.yaml
