#!/bin/bash

if [ -x "$(command -v curl)" ]; then
  curl --silent http://neverssl.com > 1
  curl --silent --cacert ca_mitm.crt https://www.google.com > 2
fi

if [ -x "$(command -v dog)" ]; then
  dog --https @https://dns.google/dns-query www.youtube.com > 3
fi
