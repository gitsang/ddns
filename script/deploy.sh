#!/bin/bash

VERSION=v0.0.10
BINPATH=/usr/local/bin
SERVICEPATH=/etc/systemd/system

wget https://github.com/gitsang/ddns/releases/download/${VERSION}/ddns -O ${BINPATH}/ddns
wget https://github.com/gitsang/ddns/releases/download/${VERSION}/ddns.service -O ${SERVICEPATH}/ddns.service

systemctl enable ddns.service
systemctl restart ddns.service
