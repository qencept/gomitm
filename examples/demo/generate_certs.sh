#!/bin/bash
openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp256r1) -days 365 \
        -subj "/CN=Root CA for MITM" -keyout ca_mitm.key -out ca_mitm.crt
