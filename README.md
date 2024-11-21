# ddns

ddns service

## run docker

```sh
mkdir -p /data/ddns/conf /data/ddns/log
docker rm -f ddns
docker run -d \
    --name ddns \
    --restart always \
    --network host \
    -v ./conf:/opt/ddns/conf \
    -v ./log:/opt/ddns/log \
    gitsang/ddns:latest
``` 

## config

example of `./conf/ddns.yml`

```yml
accesskeyid: <aliyun_ram_account_accesskeyid>
accesskeysecret: <aliyun_ram_account_accesskeysecret>
domain: <main_domain_name>
updateintervalmin: 60
ddnss:
  - enable: false
    type: "A"
    rr: "template.home"
    interface: "ens18"
    prefix: "192.168"
  - enable: false
    type: "A"
    rr: "*.template.home"
    interface: "ens18"
    prefix: "192.168"
  - enable: false
    type: "AAAA"
    rr: "template6.home"
    interface: "ens18"
    prefix: "240e"
  - enable: false
    type: "AAAA"
    rr: "*.template6.home"
    interface: "ens18"
    prefix: "240e"
```
