#!/bin/bash

if [ ! -f ca_mitm.crt ] || [ ! -f ca_mitm.key ]; then
    ./generate_scripts.sh
fi

sudo sysctl -w net.inet.ip.forwarding=1
sudo pfctl -ef pf.conf

go build -o gomitm ../../cmd/gomitm
sudo ./gomitm --config config.yaml
