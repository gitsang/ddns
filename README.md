# DDNS

ddns service

## Run with Compose

Docker image from [https://hub.docker.com/r/gitsang/ddns](https://hub.docker.com/r/gitsang/ddns)

```
services:
  ddns:
    container_name: ddns
    image: gitsang/ddns:latest
    restart: always
    network_mode: host
    volumes:
      - ./config.yaml:/etc/ddns/config.yaml
      - /var/log/ddns:/var/log/ddns
    command: -c /etc/ddns/config.yaml
```

Config file example see [./configs/example.yaml](./configs/example.yaml)
