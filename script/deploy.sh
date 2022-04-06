#!/bin/bash

VERSION=v0.0.10
BIN_PATH=/usr/local/bin
SERVICE_PATH=/etc/systemd/system
CONF_PATH=/usr/local/etc/ddns

wget https://github.com/gitsang/ddns/releases/download/${VERSION}/ddns -O ${BIN_PATH}/ddns
wget https://github.com/gitsang/ddns/releases/download/${VERSION}/ddns.service -O ${SERVICE_PATH}/ddns.service
wget https://github.com/gitsang/ddns/releases/download/${VERSION}/ddns.yml -O ${CONF_PATH}/ddns.yml

systemctl enable ddns.service
systemctl restart ddns.service
